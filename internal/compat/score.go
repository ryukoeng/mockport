package compat

type Score struct {
	Adapter          string `json:"adapter"`
	Level            string `json:"level"`
	Total            int    `json:"total"`
	EndpointCoverage int    `json:"endpoint_coverage"`
	ScenarioCoverage int    `json:"scenario_coverage"`
	SDKCoverage      int    `json:"sdk_coverage"`
	StateCoverage    int    `json:"state_coverage"`
	ErrorCoverage    int    `json:"error_coverage"`
}

func CalculateScore(manifest Manifest) Score {
	score := Score{
		Adapter:          manifest.Adapter,
		Level:            string(highestLevel(manifest.Levels)),
		EndpointCoverage: endpointCoverage(manifest.Endpoints),
		ScenarioCoverage: scenarioCoverage(manifest.Scenarios),
	}
	if hasLevel(manifest.Levels, LevelSDK) && len(manifest.SDKVersions) > 0 {
		score.SDKCoverage = 100
	}
	if hasLevel(manifest.Levels, LevelClient) && len(manifest.ClientEvidence) > 0 {
		score.SDKCoverage = 100
	}
	if hasLevel(manifest.Levels, LevelState) {
		score.StateCoverage = 100
	}
	if hasLevel(manifest.Levels, LevelError) {
		score.ErrorCoverage = 100
	}
	score.Total = (score.EndpointCoverage + score.ScenarioCoverage + score.SDKCoverage + score.StateCoverage + score.ErrorCoverage) / 5
	return score
}

func CanPromote(manifest Manifest, score Score, target string) bool {
	switch target {
	case "experimental":
		return true
	case "sdk-compatible":
		hasSDKOrClientEvidence := hasLevel(manifest.Levels, LevelSDK) || hasLevel(manifest.Levels, LevelClient)
		return hasSDKOrClientEvidence && score.SDKCoverage == 100 && score.Total >= 40
	case "workflow-compatible":
		return hasLevel(manifest.Levels, LevelWorkflow) && hasLevel(manifest.Levels, LevelState) && hasLevel(manifest.Levels, LevelError) && score.Total >= 60
	case "provider-compatible":
		return hasLevel(manifest.Levels, LevelContract) && score.Total >= 80
	default:
		return false
	}
}

func endpointCoverage(endpoints []Endpoint) int {
	if len(endpoints) == 0 {
		return 0
	}
	supported := 0
	for _, endpoint := range endpoints {
		if endpoint.Supported {
			supported++
		}
	}
	return supported * 100 / len(endpoints)
}

func scenarioCoverage(scenarios []Scenario) int {
	total := 0
	supported := 0
	for _, scenario := range scenarios {
		if !scenario.BuiltIn {
			continue
		}
		total++
		if scenario.Supported {
			supported++
		}
	}
	if total == 0 {
		return 0
	}
	return supported * 100 / total
}

func highestLevel(levels []Level) Level {
	for _, level := range []Level{LevelContract, LevelWorkflow, LevelState, LevelSDK, LevelWire, LevelError} {
		if hasLevel(levels, level) {
			return level
		}
	}
	return LevelWire
}

func hasLevel(levels []Level, want Level) bool {
	for _, level := range levels {
		if level == want {
			return true
		}
	}
	return false
}
