package management

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"

	"net/http"

	"github.com/Ceald1/purrimeter/api/db"
)

type Auth struct {
	JwtKey string `json:"jwtKey" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IiJ9.TrozRjDs4mRJ3yh9QMexo3yVJVTmOr8MTAkbVsFSudA"`

}


type AgentConfig struct {
	Name string `yaml:"name"`
	Groups []string `yaml:"groups"`
	Channels []string `yaml:"channels"`
}

// @BasePath /api/v1

// Get Agent configuration
// @Summary get agent configuration
// @Schemes
// @Description Get agent configuration
// @Tags agent management
// @Accept json
// @Produce plain
// @Success 200 {AgentConfig} agent configuration
// @Router /agent/management/getConfig [get]
func GetAgentConfig( g *gin.Context, sqlite *sql.DB, valkey valkey.Client) {
	// var auth Auth
	// err := g.ShouldBindBodyWithJSON(&auth)
	// if err != nil {
	// 	g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	// 	return
	// }
	var auth Auth
	auth.JwtKey = g.GetHeader("Authorization")
	jwtClaims, err := db.DecodeToken(auth.JwtKey)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	agentName := jwtClaims.Username
	err = db.Valkey_FetchAgent(valkey, agentName)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	agent_conf, err := db.Get_Agent_config(agentName, sqlite)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	var agent_data = AgentConfig{
		Name: agentName,
		Groups: agent_conf["groups"].([]string),
		Channels: agent_conf["channels"].([]string),
	}
	g.YAML(http.StatusOK, agent_data)

}


// @BasePath /api/v1

// Update Agent configuration
// @Summary Update agent configuration
// @Schemes
// @Description Update agent configuration
// @Tags agent management
// Param agent_config body AgentConfig true "configuration update"
// @Accept plain
// @Produce json
// @Success 200 {string} status
// @Router /agent/management/updateConfig [post]
func UpdateAgentConfig( g *gin.Context, sqlite *sql.DB, valkey valkey.Client) {
	var agent_conf AgentConfig
	auth := g.GetHeader("Authorization")
	if len(auth) < 10 {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Errorf("authorization token required!").Error()})
		return
	}
	jwtClaims, err := db.DecodeToken(auth)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	agentName := jwtClaims.Username
	err = db.Valkey_FetchAgent(valkey, agentName)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	err = g.BindYAML(agent_conf)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	err = db.Update_Agent_config(agentName, sqlite, agent_conf)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	g.JSON(200, "ok")

}