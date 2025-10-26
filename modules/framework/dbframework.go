package framework

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	manticoresearch "github.com/manticoresoftware/manticoresearch-go"
)


var ctx = context.Background()



type DBFramework struct {
	MantiClient *manticoresearch.APIClient
	// handler func(map[string]interface{}) (map[string]interface{}, error)
}


// func NewDBModulehandler (handler func(map[string]interface{}) (map[string]interface{}, error), DBHost string) *DBFramework {
// 	newconfig := manticoresearch.NewConfiguration()
// 	newconfig.Servers[0].URL = DBHost
// 	client := manticoresearch.NewAPIClient(newconfig)
// 	db := &DBFramework{
// 		MantiClient: client,
// 		handler: handler,
// 	}
// 	return db
// }
func NewDBModulehandler (DBHost string) *DBFramework {
	newconfig := manticoresearch.NewConfiguration()
	newconfig.Servers[0].URL = DBHost
	client := manticoresearch.NewAPIClient(newconfig)
	db := &DBFramework{
		MantiClient: client,
	}
	return db
}


func ( s *DBFramework) SQL(query string) (result map[string]interface{}, err error) {
	// full SQL querying
	_, httpResp, err := s.MantiClient.UtilsAPI.Sql(ctx).Body(query).Execute()
	if err != nil {
		return
	}
	if httpResp == nil {
		return
	}
	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	return
}


func (s *DBFramework) MatchJSON(match_value string, matchKey string, indexName string, limit *int32) (result map[string]interface{}, err error) {
	// searches using JSON
	SearchRequest := manticoresearch.NewSearchRequest(indexName)
	searchQuery := manticoresearch.NewSearchQuery() 
	searchQuery.Match = map[string]interface{}{
		matchKey: match_value,
	}
	SearchRequest.Query = searchQuery
	SearchRequest.Limit = limit

	_, httpResp, err := s.MantiClient.SearchAPI.Search(ctx).SearchRequest(*SearchRequest).Execute()
	if err != nil {
		return
	}
	if httpResp == nil {
		return
	}
	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	return
}



func (s *DBFramework) MatchLargerThanJSON(match_value string, matchKey string, largerThan any, limit *int32, indexName string) (result map[string]interface{}, err error) {
	SearchRequest := manticoresearch.NewSearchRequest(indexName)
	searchQuery := manticoresearch.NewSearchQuery()
	searchQuery.Range = map[string]interface{}{
		matchKey: map[string]interface{}{
			match_value: largerThan,
		},
	}
	SearchRequest.Query = searchQuery
	SearchRequest.Limit = limit

	_, httpResp, err := s.MantiClient.SearchAPI.Search(ctx).SearchRequest(*SearchRequest).Execute()
	if err != nil {
		return
	}

	if httpResp == nil {
		return
	}

	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	return
}


// copied from DB package
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

func manti_SubmitLogRaw(client *manticoresearch.APIClient, log map[string]interface{}, index string) error {  
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
  
    insertReq := manticoresearch.NewInsertDocumentRequest(index, formatted_log)  
    insertReq.SetId(generateLogID(string(stringifiedLog), timestamp))  
  
    _, _, err = client.IndexAPI.Insert(ctx).InsertDocumentRequest(*insertReq).Execute()  
    return err  
}

func (s *DBFramework) AddLog(log map[string]interface{}, index string) (err error) {
	err = manti_SubmitLogRaw(s.MantiClient, log, index)
	return err
}

