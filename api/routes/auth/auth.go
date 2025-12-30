package auth

// handles ALL authentication and database related actions
import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ceald1/purrimeter/api/crypto"

	"encoding/base64"

	"github.com/gin-gonic/gin"
	YAML "github.com/goccy/go-yaml"
	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

var ctx = context.Background()

// register a new agent
func RegisterAgent(c *gin.Context, db *surrealdb.DB) {
	var agent_secret = c.GetHeader("Authorization") // secret key for registering new agents
	agent_secret = strings.Replace(agent_secret, "Bearer ", "", -1)

	challenge := checkAgentSecret(db, agent_secret)
	if challenge != nil {
		c.JSON(403, ErrorResponse{Error: challenge.Error()}) // you fr??
		return
	}
	var agentRegister AgentRegister
	err := c.ShouldBindBodyWithJSON(&agentRegister)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	b64Decoded, err := base64.RawStdEncoding.DecodeString(agentRegister.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})  // didn't base64 encode it
		return
	}

	// check configuration file before sending to database.
	var test AgentConfig
	err = YAML.Unmarshal(b64Decoded, &test)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()}) // dumbass didn't read the YAML configuration docs
		return
	}

	// update database
	err = CreateOrUpdateAgent(db, agentRegister.Name, b64Decoded)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()}) // uh oh
		return
	}
	agentToken, err := crypto.CreateToken(agentRegister.Name,false) // create jwt
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()}) // uh oh 2
	}

	c.JSON(http.StatusOK, Result{Result: agentToken}) // finally return jwt if all is gud 😄
}















// --- Database operations ----

// check agent secret, returns error if the challenge fails or secret provided is invalid
func checkAgentSecret(db *surrealdb.DB, secret_challenge string) (err error) {
	secret_challenge = crypto.Hash(secret_challenge)
	err = db.Use(ctx, `secrets`, `secrets`)
	secretID := models.NewRecordID(`secrets`, secret_challenge)
	secret, err := surrealdb.Select[AgentSecret](ctx, db, secretID)
	if err != nil {
		return err
	}
	if secret == nil {
		return fmt.Errorf("invalid secret")
	}
	return nil
}






// create a new agent or update agent entry with the agent name and YAML configuration.
func CreateOrUpdateAgent(db *surrealdb.DB, agent_name string, config []byte) (err error) {
	err = db.Use(ctx, `agents`, `agents`)
	agent_data := Agent{
		Name: agent_name,
		Config: config, // byte array containing YAML configuration
	}
	recordID := models.NewRecordID(`agents`, agent_name)
	agent, err := surrealdb.Select[Agent](ctx, db, recordID)
	if err == nil && agent == nil {
		_, err = surrealdb.Create[Agent](ctx, db, recordID, agent_data)
		return  err
	}else{
		if agent != nil {
			_, err = surrealdb.Update[Agent](ctx, db, recordID, agent_data)
			return err
		}
	}
	return err
}



// create user table and data if doesn't exist and input agent secrets into database
func INIT(db *surrealdb.DB, api_user, api_user_pass, agent_secret string) (err error) {
	agent_secret = crypto.Hash(agent_secret)
	err = db.Use(ctx, `users`, `users`)
	if err != nil {
		return err
	}
	user_data := User{
		Name: api_user,
		Password: crypto.Hash(api_user_pass),
	}
	recordID := models.NewRecordID("users", api_user)
	user, err := surrealdb.Select[User](ctx, db, recordID)
	
	if err == nil && user == nil {
		_, err = surrealdb.Create[User](ctx, db, recordID, user_data)
		return err
	}
	if err != nil {
		return err
	}


	// agent secret
	err = db.Use(ctx, `secrets`, `secrets`)
	if err != nil {
		return  err
	}
	secretID := models.NewRecordID(`secrets`, agent_secret)
	agent_secret_data := AgentSecret{
		Name: agent_secret,
	}
	secret, err := surrealdb.Select[AgentSecret](ctx, db, secretID)
	if err != nil {
		return err
	}
	if secret == nil {
		_, err = surrealdb.Create[AgentSecret](ctx, db, secretID, agent_secret_data)
	}
	return err
}

