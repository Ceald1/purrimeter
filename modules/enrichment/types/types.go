package enrichment

type Pipeline struct {
    Pipeline PipelineStep `yaml:"pipeline"`
}


type PipelineStep map[string]Step

type Step struct {
    Namespace string   `yaml:"namespace"`
    Database  string   `yaml:"database"`
    Table     string   `yaml:"table"`
    Fields    []string `yaml:"fields"`
	Query string `yaml:"query,omitempty"`
	PushTO string `yaml:"pushTO"`
}

// --- API shit ---


// --- Default API typings ---
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

