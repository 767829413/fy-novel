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
type Info struct {
	Base struct {
		SourceID     int    `mapstructure:"source-id" json:"source-id"`
		DownloadPath string `mapstructure:"download-path" json:"download-path"`
		Extname      string `mapstructure:"extname" json:"extname"`
		LogLevel     string `mapstructure:"log-level" json:"log-level"`
	} `mapstructure:"base"  json:"base"`
	Crawl struct {
		Threads int `mapstructure:"threads" json:"threads"`
	} `mapstructure:"crawl" json:"crawl"`
	Retry struct {
		MaxAttempts int `mapstructure:"max-attempts" json:"max-attempts"`
	} `mapstructure:"retry" json:"retry"`
}

func init() {
	confValue.Store(Info{})
	loadConfig()
}

// ToJSON returns the JSON string representation of the Config
func (i Info) ToJSON() (string, error) {
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

	var newConf Info
	if err := viper.Unmarshal(&newConf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	confValue.Store(newConf)
	return nil
}

func GetConf() Info {
	return confValue.Load().(Info)
}
