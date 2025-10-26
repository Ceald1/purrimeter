package framework
import (
	"github.com/gin-gonic/gin"
)


type ModuleHTTPFramework struct {
	engine *gin.Engine
	handler func(map[string]interface{}) (map[string]interface{}, error)
}

func (s *ModuleHTTPFramework) handle(c *gin.Context) {
	var input map[string]interface{}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	output, err := s.handler(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, output)
}



func NewModuleHTTP(handler func(map[string]interface{}) (map[string]interface{}, error)) *ModuleHTTPFramework {
	gin.SetMode(gin.ReleaseMode)
	service := &ModuleHTTPFramework{
		engine:  gin.New(),
		handler: handler,
	}
	
	service.engine.POST("/", service.handle)
	
	return service
}

func (s *ModuleHTTPFramework) RunModuleHTTP(port string) error {
	return s.engine.Run(":" + port)
}
