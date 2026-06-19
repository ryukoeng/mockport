package adapter

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ScenarioHeader はリクエスト単位のシナリオ切り替えに使うヘッダ名。
const ScenarioHeader = "X-Mockport-Scenario"

// ErrUnknownScenario は未知のシナリオ名が指定されたときに返るエラー。
var ErrUnknownScenario = errors.New("unknown scenario")

// ScenarioResolver はリクエスト単位のシナリオ解決ロジックを担う。
// 解決順序: ヘッダ > Config.Scenario > defaultName。
type ScenarioResolver struct {
	configured  string
	defaultName string
	known       map[string]bool
}

// NewScenarioResolver は ScenarioResolver を生成する。
func NewScenarioResolver(cfg Config, defaultName string, meta Metadata) *ScenarioResolver {
	known := make(map[string]bool, len(meta.Scenarios))
	for _, s := range meta.Scenarios {
		known[s.Name] = true
	}
	return &ScenarioResolver{
		configured:  cfg.Scenario,
		defaultName: defaultName,
		known:       known,
	}
}

// Resolve はリクエストのヘッダを参照してシナリオ名を解決する。
// ヘッダに未知の値が含まれる場合は ErrUnknownScenario を返す。
func (s *ScenarioResolver) Resolve(req *http.Request) (string, error) {
	if value := strings.TrimSpace(req.Header.Get(ScenarioHeader)); value != "" {
		if !s.known[value] {
			return "", fmt.Errorf("%w: %s", ErrUnknownScenario, value)
		}
		return value, nil
	}
	if s.configured != "" {
		return s.configured, nil
	}
	return s.defaultName, nil
}
