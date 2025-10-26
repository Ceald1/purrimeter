package logs

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	openapi "github.com/manticoresoftware/manticoresearch-go"
	"github.com/valkey-io/valkey-go"

	"net/http"

	"os"

	"github.com/Ceald1/purrimeter/api/db"
)

type SysLog struct {
    RawJsonLog map[string]interface{} `json:"raw_log" example:"{\"timestamp\":\"2025-10-26 01:53:27\",\"logger\":\"Ceald1\",\"data\":{\"something\":\"aaa\",\"id\":12345}}"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`

}

type Result struct {
	Result string `json:"result" example:"ok"`
}



// @BasePath /api/v1

// Send logs
// @Summary send agent logs
// @Schemes
// @Description register agent
// @Tags agent management
// @Param sysLog body		SysLog true	"Send agent logs"
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /agent/logs/agentLogs [post]
func AgentLogs(valkey valkey.Client, g *gin.Context, manticore *openapi.APIClient) {
	debug := os.Getenv("DEBUG")
	if debug != "" {
		fmt.Println("=== NEW REQUEST ===")
		fmt.Printf("Method: %s\n", g.Request.Method)
		fmt.Printf("Content-Type: %s\n", g.GetHeader("Content-Type"))
		fmt.Printf("Content-Length: %d\n", g.Request.ContentLength)
		fmt.Printf("Authorization: %s\n", g.GetHeader("Authorization"))
	}
	
	// valkey, err := db.Valkey_Init()
	// if err != nil {
	// 	fmt.Printf("Valkey init error: %v\n", err)
	// 	g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	// 	return
	// }
	
	// manticore := db.Manti_Init()
	
	jwtToken := g.GetHeader("Authorization")
	jwt_claims, err := db.DecodeToken(jwtToken)
	if err != nil {
		fmt.Printf("JWT decode error: %v\n", err)
		g.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}
	
	agent_name := jwt_claims.Username
	if debug != "" {
		fmt.Printf("Agent name: %s\n", agent_name)
	}
	
	err = db.Valkey_FetchAgent(valkey, agent_name)
	if err != nil {
		if debug != "" {
			fmt.Printf("Valkey fetch error: %v\n", err)
		}
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	// Read raw body for debugging
	bodyBytes, _ := io.ReadAll(g.Request.Body)
	if debug != "" {
		fmt.Printf("Raw body (%d bytes): %s\n", len(bodyBytes), string(bodyBytes))
	}
	// Restore body for binding
	g.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	
	var sysLog SysLog
	if err := g.ShouldBindJSON(&sysLog); err != nil {
		if debug != "" {
			fmt.Printf("Bind error: %v\n", err)
		}
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("bind error: %v", err)})
		return
	}
	if debug != ""{fmt.Printf("Parsed sysLog: %+v\n", sysLog)}

	if sysLog.RawJsonLog == nil || len(sysLog.RawJsonLog) == 0 {
		if debug != "" {
			fmt.Println("raw_log is empty")
		}
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: "raw_log is required"})
		return
	}
	if debug != "" {
		fmt.Printf("Submitting to Manticore...\n")
	}
	err = db.Manti_SubmitLogRaw(manticore, sysLog.RawJsonLog)
	if err != nil  && !strings.Contains(err.Error(), "409"){
		if debug != "" {
			fmt.Printf("Manticore error: %v\n", err)
		}
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if debug != "" {
		fmt.Println("Success!")
	}
	g.JSON(http.StatusOK, Result{Result: "ok"})
}