package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
)

//go:embed default_config.yaml
var defaultConfig embed.FS

var (
	confValue atomic.Value
	initOnce  sync.Once
)

// User-defined configuration paths
const customConfigPath = "$HOME/.fynovel/config.json"

// { "base": { "source-id": 3, "download-path": "downloads", "extname": "epub", "log-level": "error" }, "crawl": { "threads": -1 }, "retry": { "max-attempts": 3 } }

// Config stores all configuration of the application.
type Info struct {
	Base struct {
		SourceID     int    `mapstructure:"source-id" json:"source-id"`
		DownloadPath string `mapstructure:"download-path" json:"download-path"`
		Extname      string `mapstructure:"extname" json:"extname"`
		LogLevel     string `mapstructure:"log-level" json:"log-level"`
	} `mapstructure:"base"    json:"base"`
	Crawl struct {
		Threads int `mapstructure:"threads" json:"threads"`
	} `mapstructure:"crawl"   json:"crawl"`
	Retry struct {
		MaxAttempts int `mapstructure:"max-attempts" json:"max-attempts"`
	} `mapstructure:"retry"   json:"retry"`
	Chatbot struct {
		Model string `mapstructure:"model" json:"model"`
	} `mapstructure:"chatbot" json:"chatbot"`
}

func init() {
	initOnce.Do(func() {
		confValue.Store(Info{})
		if err := loadConfig(); err != nil {
			// In the event of an initialization failure, we should log the error or take appropriate action
			fmt.Printf("Failed to load config: %v\n", err)
		}
	})
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

	// Check for the existence of a user-defined configuration file
	expandedPath := os.ExpandEnv(customConfigPath)
	if _, err := os.Stat(expandedPath); err == nil {
		// User-defined configuration file exists, load it
		viper.SetConfigFile(expandedPath)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read custom config file: %w", err)
		}
	} else {
		// User-defined configuration file does not exist, load default configuration
		defaultConfigFile, err := defaultConfig.Open("default_config.yaml")
		if err != nil {
			return fmt.Errorf("failed to open default config: %w", err)
		}
		defer defaultConfigFile.Close()

		viper.SetConfigType("yaml")
		if err := viper.ReadConfig(defaultConfigFile); err != nil {
			return fmt.Errorf("failed to read default config: %w", err)
		}
	}

	// Allow overriding configuration via environment variables
	viper.AutomaticEnv()

	// Parsing Configuration into Structures
	var newConf Info
	if err := viper.Unmarshal(&newConf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Storing new configurations
	confValue.Store(newConf)
	return nil
}

func GetConf() Info {
	return confValue.Load().(Info)
}

func SetConf(conf string) error {
	var newConf Info
	currentConf := GetConf()

	// Parse the input JSON string
	if err := json.Unmarshal([]byte(conf), &newConf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Compare and update configuration
	updated := false

	// Update Base fields
	if newConf.Base.SourceID != 0 && newConf.Base.SourceID != currentConf.Base.SourceID {
		currentConf.Base.SourceID = newConf.Base.SourceID
		updated = true
	}
	if newConf.Base.DownloadPath != "" &&
		newConf.Base.DownloadPath != currentConf.Base.DownloadPath {
		currentConf.Base.DownloadPath = newConf.Base.DownloadPath
		updated = true
	}
	if newConf.Base.Extname != "" && newConf.Base.Extname != currentConf.Base.Extname {
		currentConf.Base.Extname = newConf.Base.Extname
		updated = true
	}
	if newConf.Base.LogLevel != "" && newConf.Base.LogLevel != currentConf.Base.LogLevel {
		currentConf.Base.LogLevel = newConf.Base.LogLevel
		updated = true
	}

	// Update Crawl fields
	if newConf.Crawl.Threads != 0 && newConf.Crawl.Threads != currentConf.Crawl.Threads {
		currentConf.Crawl.Threads = newConf.Crawl.Threads
		updated = true
	}

	// Update Retry fields
	if newConf.Retry.MaxAttempts != 0 &&
		newConf.Retry.MaxAttempts != currentConf.Retry.MaxAttempts {
		currentConf.Retry.MaxAttempts = newConf.Retry.MaxAttempts
		updated = true
	}

	// Update Chatbot fields
	if newConf.Chatbot.Model != "" && newConf.Chatbot.Model != currentConf.Chatbot.Model {
		currentConf.Chatbot.Model = newConf.Chatbot.Model
		updated = true
	}

	// If no updates, return early
	if !updated {
		return nil
	}

	// Update the in-memory configuration
	confValue.Store(currentConf)

	// Write the new configuration to the user-defined configuration file
	expandedPath := os.ExpandEnv(customConfigPath)

	// Ensure the directory exists
	dir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Convert configuration to JSON
	jsonData, err := json.MarshalIndent(currentConf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(expandedPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
