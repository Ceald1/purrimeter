package auth

import "github.com/surrealdb/surrealdb.go/pkg/models"


type Agent struct { // goes into agents database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"`
	Config		[]byte `json:"agent_config"` // yaml agent config
}

type AgentLog struct { // goes into log database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[interface{}]interface{} `json:"log_data"` // log data

}

type User struct { // for managing the API like agents, goes into users database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"`
	Password	string `json:"password"` // hashed password
}

type AgentSecret struct { // agent secrets for enrolling agents, goes into agents database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // key for enrolling agent
}