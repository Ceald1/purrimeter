package logging
import (
	"github.com/surrealdb/surrealdb.go/pkg/models"
)


// ---- DB type ----
type AgentLog struct { // goes into agentLogs database and namespace
	ID      	*models.RecordID `json:"id,omitempty"`
	Name    	string `json:"name"` // unique hash of log (this prevents duplicates)
	LogData		map[string]interface{} `json:"log_data"` // log data
	Number      int64 `json:"log_number"` // hash unique number (for alerts)

}

// --- unique API typings ----


type LogCommit struct {
	AgentName string
	LogData map[string]interface{}
}

// --- Default API typings ---
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}
type Result struct {
	Result string `json:"result" example:"ok"`
}