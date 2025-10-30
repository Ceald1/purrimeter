package users
import (
	"database/sql"


	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/Ceald1/purrimeter/api/db"
)


type LoginData struct {
	Username string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"password1234"`
}

type TokenCheck struct {
	JWTToken string `json:"jwtToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IiJ9.TrozRjDs4mRJ3yh9QMexo3yVJVTmOr8MTAkbVsFSudA"`
}


// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

type Result struct {
	Result string `json:"result" example:"ok"`
}

// @BasePath /api/v1

// User authentication
// @Summary user auth to API
// @Schemes
// @Description user auth to API
// @Tags user
// @Param loginData body		LoginData true	"Login data information"
// @Accept json
// @Produce json
// @Success 200 {string} jwt token
// @Router /users/login [post]
func Login( g *gin.Context, sqlite *sql.DB) {
	var loginData LoginData
	err := g.ShouldBindBodyWithJSON(&loginData)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err = db.AuthUser(sqlite, loginData.Username, loginData.Password)
	if err != nil {
		g.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	jwtToken, err := db.CreateToken(loginData.Username)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	g.JSON(200, jwtToken)

}


// @BasePath /api/v1

// JWT Check
// @Summary check JWT token
// @Schemes
// @Description check JWT token this is going to be used in the frontend, for agent management tokens will be more heavily checked
// @Tags user
// @Param jwtToken body		TokenCheck true	"user jwt"
// @Accept json
// @Produce json
// @Success 200 {string} ""
// @Router /users/check [post]
func Auth( g *gin.Context) {
	var jwtToken TokenCheck
	err := g.ShouldBindBodyWithJSON(&jwtToken)
	if err != nil {
		g.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	_, err = db.DecodeToken(jwtToken.JWTToken)
	if err != nil {
		g.JSON(403, ErrorResponse{Error: err.Error()})
		return
	}
	g.JSON(http.StatusOK, "")
}
