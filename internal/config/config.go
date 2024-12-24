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

// 用户自定义的配置路径
const customConfigPath = "$HOME/.fynovel/config.json"

// { "base": { "source-id": 3, "download-path": "downloads", "extname": "epub", "log-level": "error" }, "crawl": { "threads": -1 }, "retry": { "max-attempts": 3 } }

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
	initOnce.Do(func() {
		confValue.Store(Info{})
		if err := loadConfig(); err != nil {
			// 在初始化失败时，我们应该记录错误或采取适当的措施
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

	// 检查用户自定义配置文件是否存在
	expandedPath := os.ExpandEnv(customConfigPath)
	if _, err := os.Stat(expandedPath); err == nil {
		// 用户自定义配置文件存在，加载它
		viper.SetConfigFile(expandedPath)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read custom config file: %w", err)
		}
	} else {
		// 用户自定义配置文件不存在，加载默认配置
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

	// 允许通过环境变量覆盖配置
	viper.AutomaticEnv()

	// 将配置解析到结构体中
	var newConf Info
	if err := viper.Unmarshal(&newConf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 存储新的配置
	confValue.Store(newConf)
	return nil
}

func GetConf() Info {
	return confValue.Load().(Info)
}

func SetConf(conf string) error {
	var newConf Info
	// 解析 JSON 字符串
	if err := json.Unmarshal([]byte(conf), &newConf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 更新内存中的配置
	confValue.Store(newConf)

	// 将新配置写入用户自定义的配置文件
	expandedPath := os.ExpandEnv(customConfigPath)

	// 确保目录存在
	dir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 将配置转换为 JSON
	jsonData, err := json.MarshalIndent(newConf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(expandedPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
