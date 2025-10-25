package main

import (
	"fmt"
	"runtime"

	p "github.com/Ceald1/purrimeter/agent/parsers"
)

// TODO: Add logic for authenticating and re-authenticating if disconnected.
// Auth using server side secret and reauth using server side key if jwt becomes invalid

func main() {
	platform := runtime.GOOS
	switch platform{
		case "linux":
			p.RsyslogTail("/var/log/messages", "http://localhost:8080")
		default:
			fmt.Println(platform)
	}
}