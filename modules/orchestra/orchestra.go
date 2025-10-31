package main

import (
	"fmt"
	"os"
	"io"

	Framework "github.com/Ceald1/purrimeter/modules/framework"
	"encoding/json"
	"net/http"
	"time"
)

// orchestra microservice/module this module is for orchestrating where alerts go. For registering modules/integrations check the conductor component.

func main(){
	conductor_host := os.Getenv("CONDUCTOR_HOST")
	port := os.Getenv("CONDUCTOR_PORT")
	if port == "" {
		port = "8000"
	}
	var baseURL = fmt.Sprintf(`http://%s:%s/`, conductor_host, port)
	_, err := http.NewRequest(`GET`, fmt.Sprintf(`%s/sync`, baseURL), nil)
	if err != nil {
		panic(err)
	}
	var NonEnrichment = make(map[string][]map[string]map[string]string)
	var Enrichment []map[string]map[string]string
	

	
	


}