package main

import (
   "github.com/gin-gonic/gin"
   docs "github.com/Ceald1/purrimeter/api/docs"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "github.com/Ceald1/purrimeter/api/management"
   "net/http"
   "github.com/Ceald1/purrimeter/api/db"
)
// example code from: https://github.com/swaggo/gin-swagger



// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
   g.JSON(http.StatusOK,"helloworld")
}
func RedirectHTTP(g *gin.Context) {
	g.Redirect(302, "/swagger/index.html")
}



func main()  {
   // stuff for initialization
   _, err := db.SQL_Init()
   if err != nil {
      panic(err)
   }
   val, err := db.Valkey_Init()
   if err != nil {
      panic(err)
   }
   err = db.Valkey_Secrets(val)
   if err != nil {
      panic(err)
   }



   r := gin.Default()
   docs.SwaggerInfo.BasePath = "/api/v1"
   v1 := r.Group("/api/v1")
   {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld",Helloworld)
      }
      agent := v1.Group("/agent")
      {
         management_agent := agent.Group("/management")
         {
            management_agent.POST("/register", management.RegisterAgent)
         }
      }

   }
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
   r.GET("/", RedirectHTTP)
   r.GET("/api", RedirectHTTP)
   r.GET("/api/v1", RedirectHTTP)
   r.GET("/api/v1/", RedirectHTTP)
   

   // agent management

   
   
   r.Run(":8080")

}