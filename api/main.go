package main

import (
	"context"
	"log"
	"os"

	"github.com/Ceald1/purrimeter/api/routes/auth"
	"github.com/Ceald1/purrimeter/api/routes/logging"
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
  AGENT_SECRET string = os.Getenv("AGENT_SECRET")
)





func main() {
  // Create a Gin router with default middleware (logger and recovery)
  r := gin.Default()
  gin.SetMode(gin.ReleaseMode)
  db, err := surrealdb.FromEndpointURLString(ctx, "ws://surrealdb:8000") // change to `surrealdb` in prod for kubernetes
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
  err = auth.INIT(db, API_ADMIN_USER, API_ADMIN_USER, AGENT_SECRET)
  if err != nil {
    panic(err)
  }
  go logging.Async(db)
  r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
  })
  // TODO: Add stuff for user login
  // TODO: add shit for updating rules
  // TODO: Add shit for updating enrichments
  // TODO: Basic Role based access.

  // agent endpoints
  r.POST("/api/v2/agent/register", func(c *gin.Context) {
      auth.RegisterAgent(c, db)
  })
  r.POST(`/api/v2/agent/log`, func(c *gin.Context) {
    logging.SubmitLog(c, db)
  })
    r.POST(`/api/v2/agent/logs`, func(c *gin.Context) {
    logging.SubmitLogs(c, db)
  })




  if err := r.Run(); err != nil {
    log.Fatalf("failed to run server: %v", err)
  }
}