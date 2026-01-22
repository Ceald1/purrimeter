package management
import (
	"context"
	"fmt"
	"net/http"
	// "os"
	"strings"

	// CRYPTO "github.com/Ceald1/purrimeter/api/crypto"

	"encoding/base64"

	"github.com/gin-gonic/gin"
	YAML "github.com/goccy/go-yaml"
	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)
var (
	ctx = context.Background()
)


// management for agents (allat jazz)

// stuff for updating rules, enrichments, agents (all things management)
func RegisterUser(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	var userRegister UserRegister
	err = c.ShouldBindBodyWithJSON(userRegister)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Use(ctx, `alerts`, `alerts`)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	query := "DEFINE USER $username ON DATABASE PASSWORD $password ROLES $role"
    _, err = surrealdb.Query[any](ctx, db, query, map[string]interface{}{
        "username": userRegister.Username,
        "password": userRegister.Password,
        "role":     "VIEWER", // or userRegister.Role
    })
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, Result{Result: `ok`})
}



// use the token for management, needs different database variable than main one
func LoginUser(c *gin.Context) {
	db, _ := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	var userLogin UserLogin
	err := c.ShouldBindBodyWithJSON(&userLogin)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	token, err := db.SignIn(ctx, &surrealdb.Auth{
		Namespace: `alerts`,
		Database: `alerts`,
		Username: userLogin.Username,
		Password: userLogin.Password,
	})
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, Result{Result: token})
}


// update a system user on the database level.
func UpdateUser(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	var actions UpdateUserJSON
	err = c.ShouldBindBodyWithJSON(&actions)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	action := actions.Action
	switch action {
		case `add`:
			err = db.Use(ctx, actions.Namespace, actions.Database)
			if err != nil {
				c.JSON(403, ErrorResponse{Error: err.Error()})
				return
			}
			query := "DEFINE USER $username ON DATABASE PASSWORD $password ROLES $role"
			_, err = surrealdb.Query[any](ctx, db, query, map[string]interface{}{
				"username": actions.Username,
				"password": actions.UserPass,
				"role":     actions.Access,
			})
			if err != nil {
				c.JSON(403, ErrorResponse{Error: err.Error()})
				return
			}

		case `update`:
			err = db.Use(ctx, actions.Namespace, actions.Database)
			if err != nil {
				c.JSON(403, ErrorResponse{Error: err.Error()})
				return
			}
			query := "DEFINE USER OVERWRITE $username ON DATABASE PASSWORD $password ROLES $role"
			_, err = surrealdb.Query[any](ctx, db, query, map[string]interface{}{
				"username": actions.Username,
				"password": actions.UserPass,
				"role":     actions.Access,
			})
			if err != nil {
				c.JSON(403, ErrorResponse{Error: err.Error()})
				return
			}
		case `remove`:
			err = db.Use(ctx, actions.Namespace, actions.Database)
			if err != nil {
				c.JSON(403, ErrorResponse{Error: err.Error()})
				return
			}
			query := "REMOVE USER $username ON DATABASE"
			_, err = surrealdb.Query[any](ctx, db, query, map[string]interface{}{
				"username": actions.Username,
			})
			if err != nil {
				c.JSON(403, ErrorResponse{Error: err.Error()})
				return
			}
	}
	c.JSON(200, Result{Result: `ok`})
}






// agent management

// delete an agent based on name using a surrealdb token.
func DeleteAgent(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	var agent AgentDel
	err = c.ShouldBindBodyWithJSON(&agent)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// auth to database
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Use(ctx, `agents`, `agents`)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	query := fmt.Sprintf(`DELETE FROM agents WHERE name = "%s"`, agent.AgentName)
	_,err = surrealdb.Query[any](ctx, db, query, map[string]any{})
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, Result{Result: `ok`})
}

// update an agent's config based on name and using surrealDB token.
func UpdateAgent(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	var agentUpdate AgentConfigUpdate
	err = c.ShouldBindBodyWithJSON(&agentUpdate)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Use(ctx, `agents`, `agents`)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	record := models.NewRecordID(`agents`, agentUpdate.Name)
	b64Decoded, err := base64.RawStdEncoding.DecodeString(agentUpdate.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})  // didn't base64 encode it
		return
	}
	_, err = surrealdb.Update[any](ctx, db, record, map[string]any{
		`config`: b64Decoded,
	})
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return	
	}
	c.JSON(200, Result{Result: "ok"})
}

// list all tables in the rules database
func ListRuleTables(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Use(ctx, `rules`, `rules`)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, map[string][]string{`result`: Tables(db)})
}

// fetch all tables
func Tables(db *surrealdb.DB) (tables []string) {
	results, err := surrealdb.Query[map[string]interface{}](ctx, db, `INFO FOR DB`, map[string]any{})
	if err != nil {
		fmt.Println(err.Error())
		return tables
	}
	results_ := *results
	tables_raw := results_[0].Result["tables"].(map[string]interface{})
	for table := range tables_raw {
		tables = append(tables, table)
	}
	return tables
}

// search database for specific rule matching 
func SearchRules(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Use(ctx, `rules`, `rules`)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	var ruleTable  searchRules
	err = c.ShouldBindBodyWithJSON(&ruleTable)
	table := ruleTable.Table
	var query string
	if table == "" {
		tables := Tables(db)
		for _, table = range tables{
			if ruleTable.MatchStr != ""{
				if query != "" {
					query = fmt.Sprintf(`%s\nUNION\nSELECT * FROM %s WHERE "%s"`, query, table, ruleTable.MatchStr)
				}else{
					query = fmt.Sprintf(`SELECT * FROM %s WHERE "%s"`, table, ruleTable.MatchStr)
				}
			}else{
				if query != ""{
					query = fmt.Sprintf(`%s\nUNION\nSELECT * FROM %s`, query, table)
				}else{
					query = fmt.Sprintf(`SELECT * FROM %s`, table)
				}
			}
		}
	}else{
		if ruleTable.MatchStr != ""{
			query = fmt.Sprintf(`SELECT * FROM %s WHERE "%s"`, table, ruleTable.MatchStr)
		}else{
			query = fmt.Sprintf(`SELECT * FROM %s`, table)
		}
	}
	result, err := surrealdb.Query[[]SurrealRule](ctx, db, query, map[string]any{})
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	result_ := *result
	raw_results := result_[0]
	surrealRules := raw_results.Result
	var rules []Rule
	for _, rule := range surrealRules {
		rules = append(rules, rule.RuleData)
	}
	c.JSON(200, map[string]any{`result`: rules})
}


func UpdateRules(c *gin.Context) {
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000")
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	token := strings.Replace(c.GetHeader(`Authorization`), `Bearer `, ``, 1)
	if len(token) < 10 {
		c.JSON(403, ErrorResponse{Error: `No token supplied!`})
		return
	}
	err = db.Authenticate(ctx, token)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.Use(ctx, `rules`, `rules`)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	var JSData UpdateRule
	err = c.ShouldBindBodyWithJSON(&JSData)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	tableName := JSData.Table
	DecodedRule, err := base64.RawStdEncoding.DecodeString(JSData.RuleData)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	var rule Rule
	err = YAML.Unmarshal(DecodedRule, &rule)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	surrealRule := SurrealRule{
		RuleData: rule,
	}
	recordID := models.NewRecordID(tableName, surrealRule)
	_, err = surrealdb.Upsert[any](ctx, db, recordID, surrealRule)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, Result{Result: `ok`})

}