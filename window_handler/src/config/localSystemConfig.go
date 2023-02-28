package config

import (
	"encoding/json"
	"os"
)

const CONFOG_PATH = "system_config"

type LocalConfigObserver struct {
	name string
}

func (l *LocalConfigObserver) UpdateAd(config interface{}) {
	if obj, ok := config.(systemConfig); ok {
		InputConfig(obj)
	}
}

func (l *LocalConfigObserver) GetName() string {
	return l.name
}

func (l *LocalConfigObserver) SetName(name string) {
	l.name = name
}

func InputConfig(config systemConfig) {
	result, _ := json.MarshalIndent(config, "", "    ")
	_ = os.WriteFile(CONFOG_PATH, result, 0644)
}

func readConfig() systemConfig {
	var config systemConfig
	bytes, _ := os.ReadFile(CONFOG_PATH)
	_ = json.Unmarshal(bytes, &config)
	return config
}

func loadConfig() {
	config := readConfig()
	if config.Version == "" {
		loadDefaultConfig()
	} else {
		SystemConfigCache.Set(readConfig())
	}
}
