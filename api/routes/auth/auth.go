package auth

// handles ALL authentication and database related actions
import (
	"context"

	"github.com/gin-gonic/gin"
	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
	"crypto/sha256"
	"encoding/hex"
)

var ctx = context.Background()

func LoginAgent(c *gin.Context, db *surrealdb.DB) {


}


// --- cryptography and JWT tokens -----

func hash(input string) (out string) {
	h := sha256.New()
	h.Write([]byte(input))
	hashed := h.Sum(nil)
	out = hex.EncodeToString(hashed)
	return out
}








// --- Database operations ----




// create user table and data if doesn't exist
func INIT(db *surrealdb.DB, api_user, api_user_pass string) (err error) {
	
	err = db.Use(ctx, `users`, `users`)
	if err != nil {
		return err
	}
	user_data := User{
		Name: api_user,
		Password: hash(api_user_pass),
	}
	recordID := models.NewRecordID("users", api_user)
	user, err := surrealdb.Select[User](ctx, db, recordID)
	
	if err == nil && user == nil {
		_, err = surrealdb.Create[User](ctx, db, recordID, user_data)
		return err
	}
	return err
}

