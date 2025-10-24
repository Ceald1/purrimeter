package main

import (
   "github.com/gin-gonic/gin"
   docs "github.com/Ceald1/purrimeter/api/docs"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
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
   r := gin.Default()
   docs.SwaggerInfo.BasePath = "/api/v1"
   v1 := r.Group("/api/v1")
   {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld",Helloworld)
      }
   }
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
   r.GET("/", RedirectHTTP)
   r.GET("/api", RedirectHTTP)
   r.GET("/api/v1", RedirectHTTP)
   r.GET("/api/v1/", RedirectHTTP)
   
   
   r.Run(":8080")

}