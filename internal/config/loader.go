package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("load config %s: %w", path, err)
	}

	// 1段目: 通常デコード（現行どおり。失敗したら今までどおりエラー）
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	// 2段目: KnownFields(true) で再デコードを試み、失敗したら
	// SafetyWarning{Category: "unknown_config_key"} を追加（エラーにしない）
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var strictCfg Config
	if err := decoder.Decode(&strictCfg); err != nil {
		cfg.SafetyWarnings = append(cfg.SafetyWarnings, SafetyWarning{
			Field:    "config",
			Category: "unknown_config_key",
			Message:  fmt.Sprintf("config contains unrecognized keys (typo?): %v", err),
		})
	}

	ApplyDefaults(&cfg)
	if err := Validate(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ApplyDefaults(cfg *Config) {
	if cfg.Version == "" {
		cfg.Version = "0.1"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "127.0.0.1"
	}
	if cfg.Mode == "" {
		cfg.Mode = "ai-safe"
	}
	if cfg.Adapters == nil {
		cfg.Adapters = map[string]AdapterConfig{}
	}
}
