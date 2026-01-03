package logging

// Deals with agent logs being sent TO database

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ceald1/purrimeter/api/crypto"
	"github.com/gin-gonic/gin"
	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
	// "time"
)

var (
	ctx = context.Background()
	LOGS_TO_COMMIT []LogCommit
)


func SubmitLog( c *gin.Context, db *surrealdb.DB) {
	secret := os.Getenv("JWT_SECRET")
	agentToken := c.GetHeader("Authorization")
	agentToken = strings.Replace(agentToken, "Bearer ", "", -1)
	agentClaims, err := crypto.VerifyToken(agentToken, secret)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()}) // you fr??
		return
	}


	var log_data map[string]interface{}
	err = c.ShouldBindBodyWithJSON(&log_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
		return
	}
	log := LogCommit{
		AgentName: agentClaims.Name,
		LogData: log_data,
	}
	LOGS_TO_COMMIT = append(LOGS_TO_COMMIT, log)

	c.JSON(200, Result{Result: `ok`}) // all gud 😃

}

func SubmitLogs( c *gin.Context, db *surrealdb.DB) {
	secret := os.Getenv("JWT_SECRET")
	agentToken := c.GetHeader("Authorization")
	agentToken = strings.Replace(agentToken, "Bearer ", "", -1)
	agentClaims, err := crypto.VerifyToken(agentToken, secret)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()}) // you fr??
		return
	}


	var log_data []map[string]interface{}
	err = c.ShouldBindBodyWithJSON(&log_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
		return
	}
	for _, l := range log_data {
		log := LogCommit{
			AgentName: agentClaims.Name,
			LogData: l,
		}
		LOGS_TO_COMMIT = append(LOGS_TO_COMMIT, log)

	}
	c.JSON(200, Result{Result: `ok`}) // all gud 😃

}

func checkAgent(db *surrealdb.DB, agentName string) (err error) {
	err = db.Use(ctx, `agents`, `agents`)
	if err != nil {
		return err
	}
	recordID := models.NewRecordID(`agents`, agentName)
	existingAgent, err := surrealdb.Select[any](ctx, db, recordID)
	if err != nil {
		return err
	}
	if existingAgent == nil {
		return fmt.Errorf(`Agent does not EXIST!`)
	}
	return nil
}

func Async(db *surrealdb.DB) {
	for {
		for len(LOGS_TO_COMMIT) > 0 {
			l := LOGS_TO_COMMIT[0]
			err := submitLogToDB(db, l.AgentName, l.LogData)
			if err != nil {
				fmt.Println(err.Error())
			}
			LOGS_TO_COMMIT = LOGS_TO_COMMIT[1:]
		}
		time.Sleep(time.Second * 2)
	}
}


func submitLogToDB(db *surrealdb.DB, agentName string, log map[string]interface{}) error {
    err := db.Use(ctx, `agentLogs`, `agentLogs`)
    if err != nil {
        return err
    }
	var retries = 0
	var retry_limit = 20
	
    log_name := crypto.Hash(fmt.Sprintf("%v", log))
    recordID := models.NewRecordID(`agentLogs`, log_name)
    log_data := AgentLog{
        Name:    log_name,
        LogData: log,
		Number: crypto.HashToNumber(log_name).Int64(), // needs to be byte array to prevent overflow
		
    }
    
    // Just try to create - handle duplicate error if it happens and retry.
	DB:
    _, err = surrealdb.Create[AgentLog](ctx, db, recordID, log_data)
    if err != nil {
        // If already exists, that's fine - log was already recorded
        if strings.Contains(err.Error(), `already exists`) {
            return nil
        }else{
			if retries < retry_limit{
				retries = retries + 1
				time.Sleep(time.Millisecond * time.Duration(retries))
				goto DB
			}
		}
        return err
    }
    
    return nil
}