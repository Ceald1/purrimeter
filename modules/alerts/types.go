package main

import "github.com/surrealdb/surrealdb.go/pkg/models"


type AgentLog struct { // goes into agentLogs database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[string]interface{} `json:"log_data"` // log data
	Number      int64 `json:"log_number"` // hash unique number (for alerts)

}

type Alert struct {
	ID      	*models.RecordID 			`json:"id,omitempty"`
	Name    	string 						`json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[string]interface{} `json:"log_data"` // alert data
	OriginalLog *models.RecordID `json:"original_log"` // ID of original log alert originates from
}


// yaml stuff
type Rule struct {
	Streams	    []string   `yaml:"streams"`
	Conditions  Conditions `yaml:"conditions"`
	Level 		int		   `yaml:"level"`
	Description	string	   `yaml:"description"`
	ID			string	   `yaml:"id"`
	Groups      []string   `yaml:"groups"`
}

type Conditions struct {
	Contains    []string          `yaml:"contains,omitempty"` // contains can be regex or NOT
	NotContains []string          `yaml:"not_contains,omitempty"` // can be regex or NOT
	Equals      []any             `yaml:"equals,omitempty"`
	NotEquals   []any             `yaml:"not_equals,omitempty"`
	LessThan    []any             `yaml:"less_than,omitempty"`
	GreaterThan []any             `yaml:"greater_than,omitempty"`
	Field		string			  `yaml:"field"`
}
type RuleFile struct {
	Rules []Rule `yaml:"rules"`// array of rule files
}