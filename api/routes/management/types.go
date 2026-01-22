package management

import (
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type AgentDel struct {
	AgentName string `json:"agent_name"`
}

type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type UpdateUserJSON struct {
	Action string `json:"action"`
	Username string `json:"username"`
	Access string `json:"access,omitempty"`
	Namespace string `json:"namespace"`
	Database string  `json:"database"`
	UserPass string  `json:"user_password,omitempty"`
}


type AgentConfigUpdate struct {
	Name string `json:"name"` // agent hostname
	Config string `json:"config"` // base64 encoded yaml configuration
}
type searchRules struct {
	Table string `json:"table,omitempty"`
	MatchStr string `json:"match,omitempty"`
}
// update rule based on filename/table name
type UpdateRule struct {
	Table string `json:"table"`
	RuleData string `json:"rule_data"` // base64 data (easier to deal with)

}



// --- api output  default types ---
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}
type Result struct {
	Result string `json:"result" example:"ok"`
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

type SurrealRule struct {
	ID *models.RecordID `json:"id,omitempty"`
	RuleData Rule `json:"rule_data,omitempty"`
}