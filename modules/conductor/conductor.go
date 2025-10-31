package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
)



func main(){
	port := os.Getenv("CONDUCTOR_PORT")
	if port == "" {
		port = "8000"
	}
	
	modules, err := ParseModules()
	if err != nil {
		panic(err)
	}
	fmt.Println("modules are ready!")
	router := gin.Default()
	router.GET(`/module/:type`, func(ctx *gin.Context) {
		module_type := ctx.Param("type")
		if module_type == "" {
			ctx.JSON(400, "no type specified")
			return
		}
		if strings.Contains(module_type, "!"){
			excluded_modules := make(map[string][]ModuleConfig)
			for k, v := range modules {
				excluded_modules[k] = v
			}
			delete(excluded_modules, strings.Replace(module_type, "!", "", 1))
			ctx.JSON(200, excluded_modules)
			return
		}
		ctx.JSON(200, modules[module_type])
	})
	router.GET(`/sync`, func(ctx *gin.Context) {
		modules, err = ParseModules()
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200, "")
	})




	router.Run(fmt.Sprintf("0.0.0.0:%s", port))
}

type ModuleConfig struct {
	Module struct {
		Url string `toml:"url"`
		Type string `toml:"type"`
		Description string `toml:"description"`
	}

}

func ParseModules() (modules map[string][]ModuleConfig, err error){
	entries, err := os.ReadDir("./modules")
    if err != nil {
        return
    }
	// json formatted: {"enrichment": []moduleConfig, "webhook":[]moduleConfig, "somethingElse":[]moduleConfig}
	modules = make(map[string][]ModuleConfig)


	for _, entry := range entries {
		file, err := os.Open(fmt.Sprintf("./modules/%s", entry.Name()))
		if err != nil {
			return modules, err
		}
		b, err := io.ReadAll(file)
		if err != nil {
			return modules, err
		}
		var moduleConf ModuleConfig
		err = toml.Unmarshal(b, &moduleConf)
		if err != nil {
			return modules, err
		}
		moduleConf.Module.Type = strings.ToLower(moduleConf.Module.Type)
		modules[moduleConf.Module.Type] = append(modules[moduleConf.Module.Type], moduleConf)

		file.Close()

	}
	return modules, err
}