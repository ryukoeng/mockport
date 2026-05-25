package report

type Snapshot struct {
	Mode           string          `json:"mode"`
	Safety         SafetySummary   `json:"safety"`
	Adapters       []AdapterStatus `json:"adapters"`
	Requests       []Request       `json:"requests"`
	SafetyWarnings []SafetyWarning `json:"safety_warnings"`
}

type SafetySummary struct {
	Mode               string `json:"mode"`
	Safe               bool   `json:"safe"`
	RealLookingSecrets int    `json:"real_looking_secrets"`
	ExternalURLs       int    `json:"external_urls"`
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
	Field    string `json:"field"`
	Category string `json:"category"`
	Message  string `json:"message"`
}
