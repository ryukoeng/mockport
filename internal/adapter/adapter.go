package adapter

import "net/http"

type Config struct {
	BasePath             string
	Scenario             string
	FakeSecret           string
	WebhookTargetURL     string
	WebhookSigningSecret string
}

type Adapter interface {
	Name() string
	Register(mux *http.ServeMux, cfg Config) error
	FakeEnv(cfg Config) map[string]string
	Metadata() Metadata
}

type Metadata struct {
	Name              string
	Maturity          Maturity
	ProviderVersion   string
	SDKVersions       []SDKVersion
	Levels            []Level
	Capabilities      []string
	Scenarios         []Scenario
	Endpoints         []Endpoint
	StatefulResources []string
	Idempotency       bool
	Reset             bool
}

type SDKVersion struct {
	Name    string
	Version string
}

type Scenario struct {
	Name      string
	Supported bool
}

type Endpoint struct {
	Method             string
	Path               string
	SupportedScenarios []string
	Notes              string
}

type Maturity string

const (
	MaturityExperimental       Maturity = "experimental"
	MaturityPartial            Maturity = "partial"
	MaturitySDKCompatible      Maturity = "sdk-compatible"
	MaturityWorkflowCompatible Maturity = "workflow-compatible"
	MaturityProviderCompatible Maturity = "provider-compatible"
)

type Level string

const (
	LevelWire     Level = "wire"
	LevelSDK      Level = "sdk"
	LevelClient   Level = "client"
	LevelWorkflow Level = "workflow"
	LevelState    Level = "state"
	LevelError    Level = "error"
	LevelContract Level = "contract"
)

func ValidateMaturity(maturity Maturity) bool {
	switch maturity {
	case MaturityExperimental, MaturityPartial, MaturitySDKCompatible, MaturityWorkflowCompatible, MaturityProviderCompatible:
		return true
	default:
		return false
	}
}
