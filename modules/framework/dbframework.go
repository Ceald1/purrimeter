package framework

import (
	"context"
	"encoding/json"

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

