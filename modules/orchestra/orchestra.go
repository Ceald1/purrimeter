package main

import (
	"fmt"
	"os"

	Framework "github.com/Ceald1/purrimeter/modules/framework"
	"encoding/json"
	"net/http"
	"time"
)

// orchestra microservice/module this module is for orchestrating where alerts go. For registering modules/integrations check the conductor component.

func main(){
	conductor_host := os.Getenv("CONDUCTOR_HOST")
	
	


}