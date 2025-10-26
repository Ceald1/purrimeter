package management

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"

	"net/http"

	"github.com/Ceald1/purrimeter/api/db"
)

type Registration struct {
	Name string `json:"name" example:"John Doe"`
	Key  string  `json:"key" example:"f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2"`
}
// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}
type Result struct {
	Result string `json:"result" example:"ok"`
}

type NewJWTRequest struct {
	JwtKey string `json:"jwtKey" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IiJ9.TrozRjDs4mRJ3yh9QMexo3yVJVTmOr8MTAkbVsFSudA"`
}



// @BasePath /api/v1

// Agent registration
// @Summary register agent to the server
// @Schemes
// @Description register agent
// @Tags agent management
// @Param registration body		Registration true	"Agent registration information"
// @Accept json
// @Produce json
// @Success 200 {string} jwt token
// @Router /agent/management/register [post]
func RegisterAgent( g *gin.Context, sqlite *sql.DB, valkey valkey.Client) {
	// sql, err := db.SQL_Init()
	// if err != nil {
	// 	g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	// 	return
	// }
	// valkey, err := db.Valkey_Init()
	// if err != nil {
	// 	g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	// 	return
	// }

	var registerRequest Registration
	err := g.ShouldBindBodyWithJSON(&registerRequest)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	registerRequest.Name = strings.ToLower(registerRequest.Name)
	
	secret, err := db.Get_Valkey_Secrets(valkey)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if secret != registerRequest.Key {
		err = fmt.Errorf("invalid registration key!")
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.RegisterAgent(sqlite, registerRequest.Name)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	err = db.ValkeyAgentToDB(valkey, registerRequest.Name)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	jw, err := db.CreateToken(registerRequest.Name)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	g.JSON(http.StatusOK, jw)

}



