package main

import (
	"fmt"
	"os"

	Framework "github.com/Ceald1/purrimeter/modules/framework"
	"github.com/valkey-io/valkey-go"
)

// orchestra microservice/module this module is for orchestrating where alerts go. For registering modules/integrations check the conductor component.

func main(){
	valkey_host := os.Getenv("VALKEY_HOST")
	if valkey_host == "" {
		valkey_host = "valkey"
	}
	pass := os.Getenv("REDIS_PASS")

	valk, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{fmt.Sprintf("%s:6379", valkey_host)}, Password: pass})
	if err != nil {
		panic(err)
	}
	


}