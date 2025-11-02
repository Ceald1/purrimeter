package db

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
)


func Manti_Init() (*Manticoresearch.APIClient){
	MANTICORE_HOST := fmt.Sprintf("http://%s:9308",os.Getenv("MANTICORE_HOST"))
	configuration := Manticoresearch.NewConfiguration()
	configuration.Servers[0].URL = MANTICORE_HOST
	apiClient := Manticoresearch.NewAPIClient(configuration)
	return apiClient
}

func Create_Indices(client *Manticoresearch.APIClient) error {  
	query := `CREATE TABLE IF NOT EXISTS purrimeter_raw (
        ts BIGINT,
        data JSON,
        raw_text TEXT,
		log_number BIGINT

    ) type='rt'`

	_, httpResp, err := client.UtilsAPI.Sql(ctx).Body(query).Execute()
	if err != nil {
		if httpResp != nil {
			body, _ := io.ReadAll(httpResp.Body)
			fmt.Printf("Error creating purrimeter_raw: %s\n", string(body))
		}
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS purrimeter_alerts (
        ts BIGINT,
        data JSON,
        raw_text TEXT,
		alert_num BIGINT
    ) type='rt'`
	  
	_, httpResp, err = client.UtilsAPI.Sql(ctx).Body(query).Execute()
	if err != nil {
		if httpResp != nil {
			body, _ := io.ReadAll(httpResp.Body)  
			fmt.Printf("Error creating purrimeter_alerts: %s\n", string(body))
		}
		return err
	}

	return nil
}

func generateLogID(data string, ts int64) uint64 {  
    h := sha256.New()  
    h.Write([]byte(data))  
    binary.Write(h, binary.BigEndian, ts)  
    return binary.BigEndian.Uint64(h.Sum(nil)[:8])  
}

func convertTimestampToInt64(timestampStr string) (int64, error) {  
    // Parse the RFC3339 timestamp  
    t, err := time.Parse(time.RFC3339Nano, timestampStr)  
    if err != nil {  
        return 0, err  
    }  
      
    // Convert to Unix timestamp (seconds since epoch)  
    return t.Unix(), nil  
}  

func Manti_SubmitLogRaw(client *Manticoresearch.APIClient, log map[string]interface{}) error {  
    // Extract timestamp
	nextID := time.Now().UnixNano()
    timestampStr, ok := log["timestamp"].(string)  
    if !ok {  
        return fmt.Errorf("timestamp is missing or not a string")  
    }  
  
    timestamp, err := convertTimestampToInt64(timestampStr)  
    if err != nil {  
        return err  
    }  
  
    // Serialize the original log for raw_text  
    stringifiedLog, err := json.Marshal(log)  
    if err != nil {  
        return err  
    }  
  
    // Create the document for Manticore  
    formatted_log := make(map[string]interface{})  
    formatted_log["ts"] = timestamp  
    formatted_log["raw_text"] = string(stringifiedLog)
    // JSON field should be the actual map, not stringified
    formatted_log["data"] = log
	formatted_log["log_number"] = nextID
  
    insertReq := Manticoresearch.NewInsertDocumentRequest("purrimeter_raw", formatted_log)  
    insertReq.SetId(generateLogID(string(stringifiedLog), timestamp))  
    _, _, err = client.IndexAPI.Insert(ctx).InsertDocumentRequest(*insertReq).Execute()
    
    retry:
        time.Sleep(100 * time.Millisecond)
         _, _, err = client.IndexAPI.Insert(ctx).InsertDocumentRequest(*insertReq).Execute()
    if err != nil {
        if strings.Contains("EOF", err.Error()) {
            goto retry
        }
    }

    return err  
}