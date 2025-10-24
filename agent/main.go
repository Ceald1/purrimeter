package main

import (
	p "github.com/Ceald1/purrimeter/agent/parsers"
)

func main() {
	p.RsyslogTail("/var/log/messages", "http://localhost:8080")
}