package report

type Snapshot struct {
	Mode           string          `json:"mode"`
	Adapters       []AdapterStatus `json:"adapters"`
	Requests       []Request       `json:"requests"`
	SafetyWarnings []SafetyWarning `json:"safety_warnings"`
}

type AdapterStatus struct {
	Name     string `json:"name"`
	BasePath string `json:"base_path"`
	Enabled  bool   `json:"enabled"`
}

type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Status int    `json:"status"`
}

type SafetyWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
