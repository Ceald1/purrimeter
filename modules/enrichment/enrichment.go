package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Ceald1/purrimeter/api/crypto"

	YAML "github.com/goccy/go-yaml"

	"github.com/gin-gonic/gin"
	surrealdb "github.com/surrealdb/surrealdb.go"
	types "github.com/Ceald1/purrimeter/modules/enrichment/types"
)
var (
  SURREAL_ADMIN string = os.Getenv("SURREAL_ADMIN")
  SURREAL_PASS string = os.Getenv("SURREAL_PASS")
  ctx = context.Background()
)


func main() {
	// Create a Gin router with default middleware (logger and recovery)
	secret := os.Getenv("ENRICHMENT_JWT") // different JWT SECRET, this one more private and different environment variable
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000") // change to `surrealdb` in prod for kubernetes
	if err != nil {
		panic(err)
	}
	authData := &surrealdb.Auth{
		Username: SURREAL_ADMIN,
		Password: SURREAL_PASS,
	} // login data
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

	buf, err := os.ReadFile(`/app/pipeline.yaml`)
	if err != nil {
		panic(err)
	}
	var pipeline types.Pipeline
	err = YAML.Unmarshal(buf, &pipeline)
	if err != nil {
		panic(err)
	}

	// requires JWT token in header
	r.POST(`/enrichment`, func(ctx *gin.Context) {
		auth := ctx.GetHeader(`Authentication`)
		if len(auth) < 10 {
			ctx.JSON(403, types.ErrorResponse{Error: `JWT required!`}) // you fr??
			return
		}
		_, err = crypto.VerifyToken(auth, secret)
		if err != nil {
			ctx.JSON(403, types.ErrorResponse{Error: err.Error()})
			return
		}
		var  log_data map[string]interface{}
		err = ctx.ShouldBindBodyWithJSON(&log_data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
			return
		}
		ctx.JSON(http.StatusOK, enrich(log_data,"" , pipeline, db))
	})
	r.POST(`/updatePipeline`, func(ctx *gin.Context) {
		auth := ctx.GetHeader(`Authentication`)
		if len(auth) < 10 {
			ctx.JSON(403, types.ErrorResponse{Error: `JWT required!`}) // you fr??
			return
		}
		_, err = crypto.VerifyToken(auth, secret)
		if err != nil {
			ctx.JSON(403, types.ErrorResponse{Error: err.Error()})
			return
		}
		var newPipeline types.Pipeline
		err = ctx.ShouldBindBodyWithYAML(&newPipeline)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
			return
		}
		pipeline = newPipeline
		var data []byte
		data, err = YAML.Marshal(&pipeline)
		if  err != nil {
			ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
			return
		}
		err = os.WriteFile(`/app/pipeline.yaml`, data, 0644)
		if  err != nil {
			ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
			return
		}
		ctx.JSON(200, `ok`)
	})



	r.Run()
	
}

// enrich logs, specific enrichment can be used
func enrich(log map[string]interface{}, specificEnrichment string, pipeline types.Pipeline, db *surrealdb.DB) (map[string]interface{}) {
	if specificEnrichment != "" {
		return process(log, specificEnrichment, pipeline, db)
	}else{
		// loop through enrichments
		for key, _ := range pipeline.Pipeline {
			log = process(log, key, pipeline, db)
		}
		return log
	}
}

// actually process log
func process(log map[string]interface{}, enrichment string, pipeline types.Pipeline, db *surrealdb.DB) (enrichedLog map[string]interface{}) {
	enrichmentStep, ok := pipeline.Pipeline[enrichment]
	if !ok {
		return log
	}
	enrichedLog = log
	nameSpace := enrichmentStep.Namespace
	table := enrichmentStep.Table
	database := enrichmentStep.Database
	fields := enrichmentStep.Fields
	err := db.Use(ctx, nameSpace, database)
	
	var pushTO = make([]map[string]interface{}, 0)
	if err != nil {
		fmt.Println(err.Error())
		return log
	}
	for _, field := range fields{
		value := findInMap(log, field)
		if value == "" || value == nil {
			continue
		}

		var QueryResult []map[string]interface{} // convert enrichment results to json
		var query string
		// SELECT * FROM agentLogs WHERE "test" LIMIT 1 example query for searches
		if enrichmentStep.Query != ""{
			query = fmt.Sprintf(enrichmentStep.Query, table, value)
		}else{
			query = fmt.Sprintf(`SELECT * FROM %s WHERE %s IN ORDER BY DESC LIMIT 1`, table, value)
		}
		go func(){
			raw_results, err := surrealdb.Query[[]any](ctx, db, query, map[string]any{})
			if err != nil {
				panic(err)
			}
			raw_result := (*raw_results)[0].Result
			b, _ := json.Marshal(raw_result)
			json.Unmarshal(b, &QueryResult)
			pushTO = append(pushTO, QueryResult...)
		}()
	}
	enrichedLog[enrichmentStep.PushTO] = pushTO
	return enrichedLog

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
