package main

// alerts microservice/module
import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	node_num := os.Getenv("NODE_NUMBER")
	total_number_of_nodes, _ := strconv.Atoi(os.Getenv("NODE_TOTAL"))
	mod, _ := strconv.Atoi(node_num)
	// get last query
	lastQueryFile, err := os.OpenFile(fmt.Sprintf("query_%s.sql", node_num), os.O_RDWR|os.O_CREATE, 0644)
	var query string
	if err != nil {
		panic(err)
	}
	defer lastQueryFile.Close()

	// Read existing query or use default
	stat, _ := lastQueryFile.Stat()
	if stat.Size() == 0 {
		fmt.Println("no last query found....")
		
		// query = fmt.Sprintf(`SELECT * FROM purrimeter_raw WHERE log_number MOD %d = %d LIMIT 420 OFFSET 0;`, int64(total_number_of_nodes), int64(mod - 1))
		query = fmt.Sprintf(`SELECT *, log_number MOD %d = %d FROM purrimeter_raw LIMIT 420 OFFSET 0;`, int64(total_number_of_nodes), int64(mod) - 1)
		// 'SELECT *, log_number MOD 2 = 0 FROM purrimeter_raw WHERE log_number > 1 LIMIT 10'
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

	manticoreHost := os.Getenv("MANTICORE_URL")
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
			if strings.Contains(err.Error(), "EOF") {
				fmt.Println("end of data, eepy time!")
				time.Sleep(5000 * time.Millisecond)
				continue
			}
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
				check, ok := condition["check"].(string)
				if !ok {
					check = ""
				}
				
				var log_check string
				switch operator{
					case "contains":
						if check != "" {
							formatted := formatLogData(log_data)
							log_check = formatted[check]
						}else{
							log_check = raw_log
						}
						for _, value := range values {
							if strings.Contains(log_check, value) == true {
								values_found = append(values_found, value)
							}
						}
					case "regex":
						if check != "" {
							formatted := formatLogData(log_data)
							log_check = formatted[check]
						}else{
							log_check = raw_log
						}
						for _, value := range values{
							re, err := regexp.Compile(value)
							if err != nil {
								fmt.Println("regex compile error:", err)
								continue
							}
							if re.MatchString(log_check) == true {
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
					if log_data["logOrigin"] == "sysmon"{
						alert["log_data"] = map[string]interface{}{
							"eventData":map[string]interface{}{
								"technique": GetSysmonDataString(log_data, 0),
								"time": GetSysmonDataString(log_data, 1),
								"processGUID": GetSysmonDataString(log_data, 2),
								"processID": GetSysmonDataString(log_data, 3),
								"image": GetSysmonDataString(log_data, 4),
								"fileVersion": GetSysmonDataString(log_data, 5),
								"description": GetSysmonDataString(log_data, 6),
								"product": GetSysmonDataString(log_data, 7),
								"company": GetSysmonDataString(log_data, 8),
								"originalFileName": GetSysmonDataString(log_data, 9),
								"commandLine":  GetSysmonDataString(log_data, 10),
								"currentDirectory": GetSysmonDataString(log_data, 11),
								"user": GetSysmonDataString(log_data, 12),
								"logonGuid": GetSysmonDataString(log_data, 13),
								"logonId": GetSysmonDataString(log_data, 14),
								"terminalSessionId": GetSysmonDataString(log_data, 15),
								"integrityLevel": GetSysmonDataString(log_data, 16),
								"hashes": GetSysmonDataString(log_data, 17),
								"parentProcessGuid": GetSysmonDataString(log_data, 18),
								"parentProcessId": GetSysmonDataString(log_data, 19),
								"parentImage": GetSysmonDataString(log_data, 20),
								"parentCommandLine": GetSysmonDataString(log_data, 21),
								"parentUser":GetSysmonDataString(log_data, 22),
							},
							"host": log_data["host"],
							"program": log_data["program"],
							"logOrigin": log_data["logOrigin"],
						}
					}else{
						msgParsed := strings.Split(log_data["message"].(string), " ; ")
						alert["log_data"] = map[string]interface{}{
							"eventData":map[string]interface{}{
								"technique": "",
								"time": log_data["timestamp"],
								"processGUID": "",
								"processID": "",
								"image": "",
								"fileVersion": "",
								"description": "",
								"product": "",
								"company": "",
								"originalFileName": "",
								"commandLine": strings.ReplaceAll(msgParsed[3], "COMMAND=", ""),
								"currentDirectory": strings.ReplaceAll(msgParsed[1], "PWD=", ""),
								"user": strings.ReplaceAll(msgParsed[2], "USER=", ""),
								"logonGuid": "",
								"logonId": "",
								"terminalSessionId": "",
								"integrityLevel": "",
								"hashes": "",
								"parentProcessGuid": "",
								"parentProcessId": "",
								"parentImage":"",
								"parentCommandLine": "",
								"parentUser":strings.Split(msgParsed[0], " : TTY=")[0],
							},
							"host": log_data["host"],
							"program": log_data["program"],
							"logOrigin": log_data["logOrigin"],
					}}
					values_found = nil
					err = framework.AddLog(alert, `purrimeter_alerts`,timestamp)
					if err != nil && !strings.Contains(err.Error(), "409") {
						fmt.Println(err.Error())
						retry:
							time.Sleep(100 * time.Millisecond)
							err = framework.AddLog(alert, `purrimeter_alerts`,timestamp)
						if err != nil {
							if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(),"http: server closed idle connection"){
								fmt.Println("EOF Detected, probably too fast! eepy for 100 milliseconds")
								goto retry
							}
						}
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
			// 'SELECT *, log_number MOD 2 = 0 FROM purrimeter_raw WHERE log_number > 1 LIMIT 10'
			query = fmt.Sprintf(`SELECT *, log_number MOD %d = %d FROM purrimeter_raw WHERE log_number > %d LIMIT 200;`, int64(total_number_of_nodes),int64(mod - 1), int64(log_number))
			
			// query = fmt.Sprintf(`SELECT * FROM purrimeter_raw WHERE MOD log_number %d = %d AND log_number > %d LIMIT 200;`, int64(total_number_of_nodes),int64(mod - 1), int64(log_number))
			// Truncate and write new query
			lastQueryFile.Truncate(0)
			lastQueryFile.Seek(0, 0)
			lastQueryFile.WriteString(query)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

	
// Helper function
func GetNestedString(m map[string]interface{}, keys ...string) string {
    current := m
    for i, key := range keys {
        if i == len(keys)-1 {
            // Last key - get the string value
            if val, ok := current[key]; ok {
                if str, ok := val.(string); ok {
                    return str
                }
            }
            return ""
        }
        
        // Navigate deeper
        if val, ok := current[key]; ok {
            if nested, ok := val.(map[string]interface{}); ok {
                current = nested
            } else {
                return ""
            }
        } else {
            return ""
        }
    }
    return ""
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
func GetSysmonDataString(log_data map[string]interface{}, index int) string {
    // Navigate to the Data array
    if message, ok := log_data["message"].(map[string]interface{}); ok {
        if event, ok := message["Event"].(map[string]interface{}); ok {
            if eventData, ok := event["EventData"].(map[string]interface{}); ok {
                if data, ok := eventData["Data"].([]interface{}); ok {
                    if len(data) > 0 {
                        if dataItem, ok := data[index].(map[string]interface{}); ok {
                            if text, ok := dataItem["#text"].(string); ok {
                                return text
                            }
                        }
                    }
                }
            }
        }
    }
    return ""
}
func formatLogData(log_data map[string]interface{}) (formatted map[string]string) {
    // Only process if it's sysmon and hasn't been formatted yet
	formatted = make(map[string]string)
    if log_data["logOrigin"] == "sysmon" {
        // Check if already formatted
        // if _, exists := log_data["eventData"]; exists {
        //     return nil // Already formatted
        // }
        
        // Extract all Sysmon data
        sysmonData := ExtractSysmonData(log_data)
        
        // Add the parsed eventData directly to log_data
        formatted = sysmonData
        // formatted["host"] = log_data["host"]
        // formatted["program"] = log_data["program"]
		return formatted
	}
    // } else if log_data["logOrigin"] == "sudo" || log_data["logOrigin"] == "su" {
    //     // Check if already formatted
    //     if _, exists := log_data["eventData"]; exists {
    //         return nil
    //     }
        
        msgParsed := strings.Split(log_data["message"].(string), " ; ")
        if len(msgParsed) < 4 {
            return nil // Not enough parts to parse
        }
        
        formatted = map[string]string{
            "commandLine":       strings.ReplaceAll(msgParsed[3], "COMMAND=", ""),
            "currentDirectory":  strings.ReplaceAll(msgParsed[1], "PWD=", ""),
            "user":              strings.ReplaceAll(msgParsed[2], "USER=", ""),
            "parentUser":        strings.Split(msgParsed[0], " : TTY=")[0],
        }
		return formatted
    }
func ExtractSysmonData(log_data map[string]interface{}) map[string]string {
    result := make(map[string]string)
    
    message, ok := log_data["message"].(map[string]interface{})
    if !ok {
        return result
    }
    
    event, ok := message["Event"].(map[string]interface{})
    if !ok {
        return result
    }
    
    eventData, ok := event["EventData"].(map[string]interface{})
    if !ok {
        return result
    }
    
    data, ok := eventData["Data"].([]interface{})
    if !ok {
        return result
    }
    
    // Loop through all Data items and extract by @Name
    for _, item := range data {
        dataItem, ok := item.(map[string]interface{})
        if !ok {
            continue
        }
        
        name, nameOk := dataItem["@Name"].(string)
        text, textOk := dataItem["#text"].(string)
        
        if nameOk && textOk {
			name = strings.ToLower(string(name[0])) + name[1:] // for consistency with the alerts sent
            result[name] = text
        }
    }
    
    return result
}