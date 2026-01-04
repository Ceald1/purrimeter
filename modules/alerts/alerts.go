package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ceald1/purrimeter/api/crypto"
	types "github.com/Ceald1/purrimeter/modules/alerts/types"

	YAML "github.com/goccy/go-yaml"

	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)
var (
	ctx = context.Background()
	SURREAL_ADMIN string = os.Getenv("SURREAL_ADMIN")
	SURREAL_PASS string = os.Getenv("SURREAL_PASS")
	LOGGER_USER string = os.Getenv("LOGGER_USER")
	LOGGER_PASS string = os.Getenv("LOGGER_PASS")
	NUM_OF_ALERT_SERVICES, _ = strconv.Atoi(os.Getenv("NUM_OF_ALERT_SERVICES")) // for clustering alerts for scaling purposes
	ALERT_SERVICE_NUMBER, _ = strconv.Atoi(os.Getenv("ALERT_SERVICE_NUMBER")) // for clustering alerts for scaling purposes
	query string
	SURREAL_HOST = "surrealdb" // change as needed.

)

// add Database user for read operations to `agentLog` and owner permissions for `alerts`
func create_logger_user(db *surrealdb.DB, user, password string) (err error) {
	err = db.Use(ctx, `agentLogs`, `agentLogs`)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`DEFINE USER IF NOT EXISTS %s ON ROOT PASSWORD "%s" ROLES VIEWER;`, user, password)  
    _, err = surrealdb.Query[any](ctx, db, query, map[string]any{})  
	if err != nil {
		return err
	}
	err = db.Use(ctx, `alerts`, `alerts`)
	if err != nil {
		return err
	}
	query = fmt.Sprintf(`DEFINE USER IF NOT EXISTS %s ON ROOT PASSWORD "%s" ROLES OWNER`, user, password)
	_, err = surrealdb.Query[any](ctx, db, query, map[string]any{}) 
    return err  
}
func get_last_query(filename string) (query string) {
	// get last ran query if none create it and write to a file.
	lastQueryFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	stat, _ := lastQueryFile.Stat()
	if stat.Size() == 0 {
		// int64(total_number_of_nodes), int64(mod)-1
		fmt.Println("no last query found....")
		// query = `SELECT *, log_number %% $num_services = $service_index FROM agentLogs LIMIT 420`
		query = fmt.Sprintf(`SELECT *, log_number %% %d = %d FROM agentLogs LIMIT 420`, int64(NUM_OF_ALERT_SERVICES), int64(ALERT_SERVICE_NUMBER) - 1)
		lastQueryFile.WriteString(query)
	}else {
		buf := make([]byte, stat.Size())
		n, err := lastQueryFile.Read(buf)
		if err != nil {
			panic(err)
		}
		query = strings.TrimSpace(string(buf[:n]))
	}


	return query
}


func grabRules() (rules []types.Rule, err error) {
	var path string = "/app/rules"
	files, err := os.ReadDir(path)
    if err != nil {
        err = fmt.Errorf("Error reading directory: %s", err)
        return
    }
	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return rules, err
		}
		var rulesData types.RuleFile
		err = YAML.Unmarshal(data, &rulesData)
		if err != nil {
			return rules, err
		}
		rules = append(rules, rulesData.Rules...)

	}
	return rules, nil
}



func main(){
	queryFile := fmt.Sprintf(`/app/alerts_%d.sql`, ALERT_SERVICE_NUMBER)
	rule_set, err := grabRules()
	if err != nil {
		panic(err)
	}
	query = get_last_query(queryFile)
	// basic crap that needs to run when starting.
	db, err := surrealdb.FromEndpointURLString(ctx, fmt.Sprintf("ws://%s:8000", SURREAL_HOST)) // change to `surrealdb` in prod
	if err != nil {
		panic(err)
	}
	var authData *surrealdb.Auth
	if SURREAL_ADMIN != "" && SURREAL_PASS != ""{
		authData = &surrealdb.Auth{
			Username: SURREAL_ADMIN,
			Password: SURREAL_PASS,
		} // login data
	}else{
		authData = &surrealdb.Auth{
			Username: LOGGER_USER,
			Password: LOGGER_PASS,
		}
	}
	token, err := db.SignIn(ctx, authData) // sign in
		if err != nil {
			panic(err)
		}
	if err = db.Authenticate(ctx, token); err != nil {
		panic(err)
	}



	defer func(token string) {
			if err := db.Invalidate(ctx); err != nil {
				panic(err)
			}
		}(token) // delete token after function ends
	if ALERT_SERVICE_NUMBER == 1 {
		err = create_logger_user(db, LOGGER_USER, LOGGER_PASS)
		if err != nil {
			panic(err)
		}
	}
	
	var lastLog types.AgentLog
	var realtimeUpdate bool = false
	START_AGAIN:
	err = db.Use(ctx, `agentLogs`, `agentLogs`) // test agentLogs
	if err != nil {
		panic(err)
	}
	lastQueryFile, _ := os.OpenFile(queryFile, os.O_RDWR|os.O_CREATE, 0644)
	lastQueryFile.WriteString(query)
	entries, err := surrealdb.Query[[]types.AgentLog](ctx, db, query, map[string]any{})
	if err != nil {
		panic(err)
	}
	for _, entry := range *entries {
		if len(entry.Result) == 0 {
			realtimeUpdate = true
			break
		}
		lastLog = entry.Result[len(entry.Result)-1]
		for _, log := range entry.Result {
			alert(rule_set, log, db)
		}
	}
	// repeat and start again.
	// example query for offset: 
	// ```go
	// SELECT * from agentLogs:033fd53d61ee9fe965df708c89801251481e693e65035f54588bf5c55b1e99b1>.. LIMIT 1 START 1
	// ```
	if realtimeUpdate == false{
		query = fmt.Sprintf(`SELECT *, log_number %% %d = %d FROM %s>.. LIMIT 200 START 1`,int64(NUM_OF_ALERT_SERVICES), int64(ALERT_SERVICE_NUMBER) - 1, lastLog.ID)
		// fmt.Sprintf(`SELECT *, log_number %% %d = %d FROM agentLogs LIMIT 420`, int64(NUM_OF_ALERT_SERVICES), int64(ALERT_SERVICE_NUMBER) - 1)
		goto START_AGAIN
	}
	fmt.Println("starting realtime updates...")
	liveDB, err := surrealdb.Live(ctx, db, `agentLogs`, false)
	notifications, err := db.LiveNotifications(liveDB.String())

	if err != nil {
		panic(err)
	}
	// live stuff goes here.
	for notification := range notifications {
		resultAny := notification.Result.(map[string]any)
		var result types.AgentLog
		jsData, _ := json.Marshal(resultAny)
		json.Unmarshal(jsData, &result)
		resultID := result.Number
		if resultID % int64(NUM_OF_ALERT_SERVICES) == (int64(ALERT_SERVICE_NUMBER) - 1) {
			alert(rule_set, result, db)
			lastLog = result
			query = fmt.Sprintf(`SELECT * FROM %s>.. LIMIT 200 START 1`, lastLog.ID)
			lastQueryFile, _ := os.OpenFile(queryFile, os.O_RDWR|os.O_CREATE, 0644)
			lastQueryFile.WriteString(query)
		}


	}
}

func alert(rule_set []types.Rule, log types.AgentLog, db *surrealdb.DB) {
	// var err error
	for _, rule := range rule_set{
		field := rule.Conditions.Field
		description := rule.Description
		id := rule.ID
		level := rule.Level
		groups := rule.Groups
		streams := rule.Streams
		entryID := log.ID
		field_value := FindField(log, field)
		if field_value == nil {
			continue
		}
		contains := rule.Conditions.Contains
		notContains := rule.Conditions.NotContains
		equals := rule.Conditions.Equals
		notEquals := rule.Conditions.NotEquals
		lessThan := rule.Conditions.LessThan
		greaterThan := rule.Conditions.GreaterThan
		var rule_formatted map[string]interface{}
		// contains
		for _, contain := range contains {
			var match bool = false
			isRegex := IsValidRegex(contain)
			if isRegex == true {
				expression, _ := regexp.Compile(contain)
				match = expression.MatchString(field_value.(string)) // trusting string
				}else{
				match = strings.Contains(field_value.(string), contain)
			}
			if match == true {
					// send alert to DB and reference the ID
					rule_formatted = map[string]interface{}{`id`:id, `level`: level, `description`: description, `groups`: groups, `streams`: streams, `field`: field}
					// if err != nil {
					// 	panic(err)
					// }
					goto StopCheck
				}
			
		}

		// not contains
		for _, notContain := range notContains {
			var match bool = false
			isRegex := IsValidRegex(notContain)
			if isRegex == true {
				expression, _ := regexp.Compile(notContain)
				match = !expression.MatchString(field_value.(string))
			}else{
				match = !strings.Contains(field_value.(string), notContain)
			}
			if match == true {
					// send alert to DB and reference the ID
					rule_formatted = map[string]interface{}{`id`:id, `level`: level, `description`: description, `groups`: groups, `streams`: streams, `field`: field}
					// err = SendAlert(rule_formatted, entryID, db)
					// if err != nil {
					// 	panic(err)
					// }
					goto StopCheck
				}
		}

		// check if equals
		for _, equal := range equals {
			if fmt.Sprintf("%v",equal) == fmt.Sprintf("%v",field_value) {
				rule_formatted = map[string]interface{}{`id`:id, `level`: level, `description`: description, `groups`: groups, `streams`: streams, `field`: field}
				// err = SendAlert(rule_formatted, entryID, db)
				// if err != nil {
				// 	panic(err)
				// }
				goto StopCheck
			}
		}

		// check if not equals
		for _, notEqual := range notEquals {
			if fmt.Sprintf("%v", notEqual) != fmt.Sprintf("%v", field_value ){
				rule_formatted = map[string]interface{}{`id`:id, `level`: level, `description`: description, `groups`: groups, `streams`: streams, `field`: field}
				// err = SendAlert(rule_formatted, entryID, db)
				// if err != nil {
				// 	panic(err)
				// }
				goto StopCheck
			}
		}

		// less than
		for _, l := range lessThan {
			if l.(int) < field_value.(int) {
				rule_formatted = map[string]interface{}{`id`:id, `level`: level, `description`: description, `groups`: groups, `streams`: streams, `field`: field}
				// err = SendAlert(rule_formatted, entryID, db)
				// if err != nil {
				// 	panic(err)
				// }
				goto StopCheck
			}
		}

		// greater than
		for _, g := range greaterThan {
			if g.(int) > field_value.(int) {
				rule_formatted = map[string]interface{}{`id`:id, `level`: level, `description`: description, `groups`: groups, `streams`: streams, `field`: field}
				// err = SendAlert(rule_formatted, entryID, db)
				// if err != nil {
				// 	panic(err)
				// }
				goto StopCheck
			}
		}


		// stop checking
			StopCheck:
				go SendAlert(rule_formatted, entryID, db)
				continue

	}
}

// send alert to DB directly and reference the original log id
func SendAlert(alertData map[string]interface{}, originalLogID *models.RecordID, db *surrealdb.DB){
	var alert types.Alert
	var retries = 0
	var retry_limit = 20
	alertData["referencedRecord"] = originalLogID.ID
	db.Use(ctx, `alerts`, `alerts`)
	recordName := crypto.Hash(fmt.Sprintf("%v", alertData))
	alert = types.Alert{
		Name: recordName,
		LogData: alertData,
		OriginalLog: originalLogID,
	}
	recordID := models.NewRecordID(`alerts`, recordName)
	DB:
	_, err := surrealdb.Create[types.Alert](ctx, db, recordID, alert)
	if err != nil {
		if strings.Contains(err.Error(), `already exists`) {
			return
		}else{
			if retries < retry_limit{
				retries = retries + 1
				time.Sleep(time.Millisecond * time.Duration(retries))
				goto DB
			}
		}
        panic(err)
	}
}

// check if is valid regex
func IsValidRegex(pattern string) bool {
    _, err := regexp.Compile(pattern)
    return err == nil
}

// find field
func FindField(log types.AgentLog, field string) (result any) {
	var logData = log.LogData

	if value, exists := logData[field]; exists {
        return value
    }
	for _, value := range logData {
		nested, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		nestedResult := findInMap(nested, field)
		if nestedResult != nil {
			return nestedResult
		}

	}
	return nil // nothing
}


func findInMap(m map[string]interface{}, field string) any {
	if value, exists := m[field]; exists {
		return  value
	}
	for _, value := range m {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			if result := findInMap(nestedMap, field); result != nil {
				return result
			}
		}
	}
	return nil
}

