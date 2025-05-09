package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	// 创建自定义日志配置
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console", // 使用控制台格式，更易读
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        zapcore.OmitKey,
			CallerKey:      zapcore.OmitKey,  // 不输出调用者信息
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  zapcore.OmitKey,  // 不输出堆栈信息
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 带颜色的日志级别
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	
	logger, err := config.Build()
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
