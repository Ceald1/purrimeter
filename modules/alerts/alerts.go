package main
import (
	"fmt"
	"context"
	surrealdb "github.com/surrealdb/surrealdb.go"
	"os"

)
var (
	ctx = context.Background()
	SURREAL_ADMIN string = os.Getenv("SURREAL_ADMIN")
	SURREAL_PASS string = os.Getenv("SURREAL_PASS")
	LOGGER_USER string = os.Getenv("LOGGER_USER")
	LOGGER_PASS string = os.Getenv("LOGGER_PASS")

)

// add Database user for read operations to `agentLog` and owner permissions for `alerts`
func create_logger_user(db *surrealdb.DB, user, password string) (err error) {
	err = db.Use(ctx, `agentLogs`, `agentLogs`)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`DEFINE USER IF NOT EXISTS %s ON ROOT PASSWORD "%s" ROLES VIEWER;`, user, password)  
    _, err = surrealdb.Query[any](ctx, db, query, map[string]any{})  
	if err != nil {
		return err
	}
	err = db.Use(ctx, `alerts`, `alerts`)
	if err != nil {
		return err
	}
	query = fmt.Sprintf(`DEFINE USER IF NOT EXISTS %s ON ROOT PASSWORD "%s" ROLES OWNER`, user, password)
	_, err = surrealdb.Query[any](ctx, db, query, map[string]any{}) 
    return err  
}



func main(){
	// basic crap that needs to run when starting.
	db, err := surrealdb.FromEndpointURLString(ctx, "ws://127.0.0.1:8000") // change to `surrealdb` in prod
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
	err = create_logger_user(db, LOGGER_USER, LOGGER_PASS)
	if err != nil {
		panic(err)
	}


	
}