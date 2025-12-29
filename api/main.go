package main

import (
	"context"
	"log"
	"os"

	"github.com/Ceald1/purrimeter/api/routes/auth"
	// "github.com/Ceald1/purrimeter/api/routes/logging"
	// "github.com/Ceald1/purrimeter/api/routes/management"
	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)
var (
  SURREAL_ADMIN string = os.Getenv("SURREAL_ADMIN")
  SURREAL_PASS string = os.Getenv("SURREAL_PASS")
  ctx = context.Background()
  API_ADMIN_USER string = os.Getenv("API_ADMIN")
  API_ADMIN_PASS string = os.Getenv("API_ADMIN_PASS")
)





func main() {
  // Create a Gin router with default middleware (logger and recovery)
  r := gin.Default()
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
  err = auth.INIT(db, API_ADMIN_USER, API_ADMIN_USER)
  if err != nil {
    panic(err)
  }




  r.POST("/api/v2/authAgent", func(c *gin.Context) {
      auth.LoginAgent(c, db)
  })

//   // Define a simple GET endpoint
//   r.GET("/ping", func(c *gin.Context) {
//     // Return JSON response
//     c.JSON(http.StatusOK, gin.H{
//       "message": "pong",
//     })
//   })

//   // Start server on port 8080 (default)
//   // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
  if err := r.Run(); err != nil {
    log.Fatalf("failed to run server: %v", err)
  }
}