package main

import (
	"fmt"
	"runtime"
	"strings"

	"os"

	p "github.com/Ceald1/purrimeter/agent/parsers"
	s "github.com/Ceald1/purrimeter/agent/senders"
)



func main() {
	jwt_file := "token.txt"
	api_url := os.Args[1]
	agentName, _ := os.Hostname()
	// agentName := "test1234"
	var jwt string

	if len(api_url) < 2 {
		panic(`API url is required! example: http://127.0.0.1`)
	}
	if _, err := os.Stat(jwt_file); err != nil {
		fmt.Println("getting token...")
		server_token := os.Getenv("AUTH_TOKEN")
		if len(server_token) < 2 {
			panic(`"AUTH_TOKEN" env variable is required for agent registration!`)
		}
		jwt, err = s.RegisterAgent(api_url, server_token, agentName)
		if err != nil {
			panic(err)
		}
		data := []byte(jwt)
		err := os.WriteFile(jwt_file, data, 0644)
		if err != nil {
			panic(err)
		}


	}else {
		raw_data, err := os.ReadFile(jwt_file)
		if err != nil {
			panic(err)
		}
		jwt = strings.ReplaceAll(string(raw_data), "\n", "")
	}
	


	platform := runtime.GOOS
	switch platform{
		case "linux":
			p.RsyslogTail("/var/log/messages",fmt.Sprintf("%s/api/v1/agent/logs/agentLogs", api_url), jwt)
		default:
			fmt.Println(platform)
	}
}