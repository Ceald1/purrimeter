package auth

import (
	"github.com/surrealdb/surrealdb.go/pkg/models"
)


type Agent struct { // goes into agents database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"`
	Config		[]byte `json:"agent_config"` // yaml agent config
}

type AgentLog struct { // goes into agentLogs database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[interface{}]interface{} `json:"log_data"` // log data

}

type User struct { // for managing the API like agents, goes into users database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"`
	Password	string `json:"password"` // hashed password
}

type AgentSecret struct { // agent secrets for enrolling agents, goes into secrets database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // key for enrolling agent
}


// ----API types----

// ---- api input types
type AgentRegister struct {
	Name string `json:"name"` // agent hostname
	Config string `json:"config"` // base64 encoded yaml configuration
}



// --- api output types ---
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}
type Result struct {
	Result string `json:"result" example:"ok"`
}

// ----yaml types----

type AgentConfig struct {
	OS 			string `yaml:"os"` // agent os
	HostName 	string `yaml:"name"` // agent hostname
	Streams		[]string `yaml:"streams"` // data streams
}