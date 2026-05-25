package report

type Snapshot struct {
	Mode                 string                `json:"mode"`
	Safety               SafetySummary         `json:"safety"`
	Adapters             []AdapterStatus       `json:"adapters"`
	Requests             []Request             `json:"requests"`
	SafetyWarnings       []SafetyWarning       `json:"safety_warnings"`
	ScenarioCoverage     []ScenarioCoverage    `json:"scenario_coverage"`
	BehaviorMatrix       []BehaviorMatrixEntry `json:"behavior_matrix"`
	UnsupportedEndpoints []UnsupportedEndpoint `json:"unsupported_endpoints"`
}

type SafetySummary struct {
	Mode               string `json:"mode"`
	Safe               bool   `json:"safe"`
	RealLookingSecrets int    `json:"real_looking_secrets"`
	ExternalURLs       int    `json:"external_urls"`
}

type AdapterStatus struct {
	Name         string   `json:"name"`
	BasePath     string   `json:"base_path"`
	Enabled      bool     `json:"enabled"`
	Scenario     string   `json:"scenario,omitempty"`
	Maturity     string   `json:"maturity,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

type Request struct {
	ID        int64  `json:"id"`
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Adapter   string `json:"adapter,omitempty"`
	Scenario  string `json:"scenario,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type SafetyWarning struct {
	Field    string `json:"field"`
	Category string `json:"category"`
	Message  string `json:"message"`
}

type ScenarioCoverage struct {
	Adapter   string            `json:"adapter"`
	Scenarios []ScenarioSupport `json:"scenarios"`
}

type ScenarioSupport struct {
	Name      string `json:"name"`
	Supported bool   `json:"supported"`
}

type BehaviorMatrixEntry struct {
	Adapter            string   `json:"adapter"`
	Maturity           string   `json:"maturity"`
	Method             string   `json:"method"`
	Path               string   `json:"path"`
	SupportedScenarios []string `json:"supported_scenarios"`
	Notes              string   `json:"notes,omitempty"`
}

type UnsupportedEndpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Status int    `json:"status"`
	Reason string `json:"reason"`
}
