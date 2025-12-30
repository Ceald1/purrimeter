package logging

// Deals with agent logs being sent TO database

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ceald1/purrimeter/api/crypto"
	"github.com/gin-gonic/gin"
	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

var (
	ctx = context.Background()
)


func SubmitLog( c *gin.Context, db *surrealdb.DB) {
	agentToken := c.GetHeader("Authorization")
	agentToken = strings.Replace(agentToken, "Bearer ", "", -1)
	agentClaims, err := crypto.VerifyToken(agentToken)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()}) // you fr??
		return
	}
	// check if agent exists
	err = checkAgent(db, agentClaims.Name)
	if err != nil {
		c.JSON(403, ErrorResponse{Error: err.Error()}) // agent no exist or error
		return
	}


	var log_data map[string]interface{}
	err = c.ShouldBindBodyWithJSON(&log_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()}) // how can you send invalid JSON??
		return
	}
	err = submitLogToDB(db, agentClaims.Name, log_data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()}) // uh oh
		return
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



func submitLogToDB(db *surrealdb.DB, agentName string, log map[string]interface{}) (err error){
	err = db.Use(ctx, `agentLogs`, `agentLogs`)
	if err != nil {
		return err
	}
	log_name := crypto.Hash(fmt.Sprintf("%v", log))
	recordID := models.NewRecordID(`agentLogs`, log_name)
	log_data := AgentLog{
		Name: log_name,
		LogData: log,
	}
	existing_log, err := surrealdb.Select[AgentLog](ctx, db, recordID)
	if err != nil {
		return err
	}
	if existing_log == nil {
		_, err = surrealdb.Create[AgentLog](ctx, db, recordID, log_data)
		return err
	}
	return nil // don't return anything if already exists just continue.
}