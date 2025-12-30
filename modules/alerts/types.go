package main

import "github.com/surrealdb/surrealdb.go/pkg/models"

type AgentLog struct { // goes into agentLogs database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[interface{}]interface{} `json:"log_data"` // log data

}

type Alert struct {
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[interface{}]interface{} `json:"log_data"` // log data
}