package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/spf13/viper"
)

//go:embed default_config.yaml
var defaultConfig embed.FS

var confValue atomic.Value

// Config stores all configuration of the application.
type info struct {
	Base struct {
		SourceID     int    `mapstructure:"source-id"`
		DownloadPath string `mapstructure:"download-path"`
		Extname      string `mapstructure:"extname"`
		AutoUpdate   int    `mapstructure:"auto-update"`
		LogLevel     string `mapstructure:"log-level"`
	} `mapstructure:"base"`
	Crawl struct {
		Threads int `mapstructure:"threads"`
	} `mapstructure:"crawl"`
	Retry struct {
		MaxAttempts int `mapstructure:"max-attempts"`
	} `mapstructure:"retry"`
}

func init() {
	confValue.Store(info{})
	loadConfig()
}

// ToJSON returns the JSON string representation of the Config
func (i info) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// LoadConfig reads configuration from file or environment variables.
func loadConfig() error {
	viper.Reset()

	// First, load the default configuration
	defaultConfigFile, err := defaultConfig.Open("default_config.yaml")
	if err != nil {
		return fmt.Errorf("failed to open default config: %w", err)
	}
	defer defaultConfigFile.Close()

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(defaultConfigFile); err != nil {
		return fmt.Errorf("failed to read default config: %w", err)
	}

	viper.AutomaticEnv()

	var newConf info
	if err := viper.Unmarshal(&newConf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	confValue.Store(newConf)
	return nil
}

func GetConf() info {
	return confValue.Load().(info)
}
