package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// ExchangeConfig 交易所配置
type ExchangeConfig struct {
	ApiKey      string `yaml:"api_key"`
	ApiSecret   string `yaml:"api_secret"`
	Passphrase  string `yaml:"passphrase,omitempty"` // OKX需要
	BaseURL     string `yaml:"base_url"`
	AccountType string `yaml:"account_type,omitempty"` // 账户类型：spot(现货)、margin(保证金)、futures(期货)，默认为spot
}

// Config 应用配置
type Config struct {
	Server struct {
		Port      string `yaml:"port"`
		Mode      string `yaml:"mode"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Exchanges map[string]ExchangeConfig `yaml:"exchanges"`
}

var AppConfig *Config
var Logger *zap.SugaredLogger

// InitConfig 初始化配置
func InitConfig() error {
	configPath := filepath.Join("configs", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}

	AppConfig = cfg
	return nil
}

// InitLogger 初始化日志
func InitLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	Logger = logger.Sugar()
	return nil
}

// GetExchangeConfig 获取指定交易所的配置
func GetExchangeConfig(name string) (*ExchangeConfig, bool) {
	if AppConfig == nil {
		return nil, false
	}
	
	cfg, exists := AppConfig.Exchanges[name]
	if !exists {
		return nil, false
	}
	
	return &cfg, true
}
