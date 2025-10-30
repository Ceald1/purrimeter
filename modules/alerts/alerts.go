package main
// alerts microservice/module
import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"
	"regexp"
	// "encoding/json"
	// "net/http"
	// "bytes"

	Framework "github.com/Ceald1/purrimeter/modules/framework"
	YAML "github.com/goccy/go-yaml"
)

var ruleList []map[string]interface{}

func isAvailable(alpha []string, str string) bool {

   // iterate using the for loop
   for i := 0;
   i < len(alpha);
   i++ {
      // check      
      if alpha[i] == str {
         // return true
         return true
      }
   }
   return false
}

func main(){
	// get last query
	lastQueryFile, err := os.OpenFile("query.sql", os.O_RDWR|os.O_CREATE, 0644)
	var query string
	if err != nil {
		panic(err)
	}
	defer lastQueryFile.Close()

	// Read existing query or use default
	stat, _ := lastQueryFile.Stat()
	if stat.Size() == 0 {
		fmt.Println("no last query found....")
		query = `SELECT * FROM purrimeter_raw LIMIT 200 OFFSET 0;`
		lastQueryFile.WriteString(query)
	} else {
		buf := make([]byte, stat.Size())
		n, err := lastQueryFile.Read(buf)
		if err != nil {
			panic(err)
		}
		query = strings.TrimSpace(string(buf[:n]))
	}
	fmt.Println(query)

	manticoreHost := os.Getenv("MANTICORE_HOST")
	if manticoreHost == "" {
		manticoreHost = "http://127.0.0.1:9308"
	}
	// orchestratorHost := os.Getenv("ORCHESTRA_HOST")
	// if orchestratorHost == "" {
	// 	orchestratorHost = "http://127.0.0.1:9080"
	// }
	framework := Framework.NewDBModulehandler(manticoreHost)
	ruleList, err := getRules()
	if err != nil {
		panic(err)
	}

	var log map[string]interface{}
	
	var log_number float64

	// search db based on unix time nano second
	for {
		queryResults, err := framework.SQL(query)
		if err != nil {
			fmt.Println(err.Error())
		}
		if queryResults == nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}
		logs_raw := queryResults["hits"].(map[string]interface{})["hits"].([]interface{})

		// Convert []interface{} to []map[string]interface{}
		logs := make([]map[string]interface{}, 0, len(logs_raw))
		for _, logInterface := range logs_raw {
			if logMap, ok := logInterface.(map[string]interface{}); ok {
				logs = append(logs, logMap)
			}
		}
		if len(logs) < 1 {
			time.Sleep(5000 * time.Millisecond)
			continue
		}
		for _, log = range logs {
			// do parsing logic with the rules
			log_data := log["_source"].(map[string]interface{})["data"].(map[string]interface{})
			raw_log := log["_source"].(map[string]interface{})["raw_text"].(string)
			source := log_data["logOrigin"].(string)
			log_number = log["_source"].(map[string]interface{})["log_number"].(float64)
			timestamp := log_data["timestamp"].(string)
			for _, rule := range ruleList {
				var rule_sources []string
				rule_sources_raw := rule["sources"].([]interface{})
				for _, r := range rule_sources_raw {
					rule_sources = append(rule_sources, r.(string))
				}
				var values_found []string
				
				condition := rule["condition"].(map[string]interface{})
				operator := condition["operator"].(string)
				values_raw := condition["values"].([]interface{})
				var values []string
				for _, value := range values_raw {
					values = append(values, value.(string))
				}
				if !isAvailable(rule_sources, source){
					continue
				}
				switch operator{
					case "contains":
						for _, value := range values {
							if strings.Contains(raw_log, value) == true {
								values_found = append(values_found, value)
							}
						}
					case "regex":
						for _, value := range values{
							re, err := regexp.Compile(value)
							if err != nil {
								fmt.Println("regex compile error:", err)
								continue
							}
							if re.MatchString(raw_log) == true {
								values_found = append(values_found, value)
							}
						}
					default:
						fmt.Println("unknown condition: ", operator)					
				}
				// craft new json data
				if len(values_found) > 0 {
					var alert = make(map[string]interface{})
					// alert["data"] = log_data
					alert["parent_log"] = log_number
					alert["alert"] = map[string]interface{}{
						"rule id": rule["id"],
						"description": rule["description"],
						"severity": rule["severity"],
						"level": rule["level"],
						"values": values_found,
						"groups": rule["groups"],
						"sources": rule["sources"],
					}
					alert["log_data"] = log_data
					values_found = nil
					err = framework.AddLog(alert, `purrimeter_alerts`,timestamp)
					if err != nil && !strings.Contains(err.Error(), "409") {
						fmt.Println(err.Error())
					}
					// else{
					// 	jsString, _ := json.Marshal(alert)
					// 	_, err = http.Post(orchestratorHost, "application/json", bytes.NewBuffer(jsString))
					// 	if err != nil {
					// 		fmt.Println(err)
					// 	}

					// }
				}
			}
			
			query = fmt.Sprintf(`SELECT * FROM purrimeter_raw WHERE log_number > %d LIMIT 200;`, int64(log_number))
			// Truncate and write new query
			lastQueryFile.Truncate(0)
			lastQueryFile.Seek(0, 0)
			lastQueryFile.WriteString(query)
		}

		time.Sleep(2000 * time.Millisecond)
	}
}

	


func getRules() (result []map[string]interface{}, err error) {
	rulesRaw, err := readAllFilesRecursive("./rules")
	if err != nil {
		return
	}
	for _, rule := range rulesRaw {
		parsed_rule, err := rule2map(rule)
		if err != nil {
			fmt.Println("error parsing rule: ", rule)
		}else{
			result = append(result, parsed_rule...)
		}
	}
	return
}


func readAllFilesRecursive(dirPath string) ([]string, error) {
    var files []string
    
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil  // Continue walking
        }
        
        content, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        
        files = append(files, string(content))
        return nil
    })
    
    return files, err
}

func rule2map(yamlStr string) (result []map[string]interface{}, err error) {
	var results map[string]interface{}
	if err = YAML.Unmarshal([]byte(yamlStr), &results); err != nil {
		return nil, err
	}
	
	rulesInterface, ok := results["rules"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("rules field not found or not an array")
	}
	
	result = make([]map[string]interface{}, len(rulesInterface))
	for i, ruleWrapper := range rulesInterface {
		// Each item has a "rule" key wrapping the actual rule
		ruleWrapperMap, ok := ruleWrapper.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("rule at index %d is not a map", i)
		}
		
		// Extract the actual rule
		actualRule, ok := ruleWrapperMap["rule"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("rule field at index %d not found or not a map", i)
		}
		
		result[i] = actualRule
	}
	
	return result, nil
}