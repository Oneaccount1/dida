package config

import (
	"os"
	"strconv"
	"time"

	"dida/internal/errors"
)

// Config 应用程序配置
type Config struct {
	// TickTick API 配置
	TickTick TickTickConfig `json:"ticktick"`

	// 服务器配置
	Server ServerConfig `json:"server"`

	// 日志配置
	Log LogConfig `json:"log"`
}

// TickTickConfig TickTick API 配置
type TickTickConfig struct {
	ClientID     string        `json:"client_id"`
	ClientSecret string        `json:"client_secret"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	BaseURL      string        `json:"base_url"`
	TokenURL     string        `json:"token_url"`
	AuthURL      string        `json:"auth_url"`
	RedirectURL  string        `json:"redirect_url"`
	Timeout      time.Duration `json:"timeout"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `json:"level"`
	FilePath string `json:"file_path"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 注意：环境变量已在 main.go 中加载，这里不需要重复加载

	config := &Config{
		TickTick: TickTickConfig{
			// 只从环境变量读取认证相关的敏感信息
			ClientID:     getEnv("TICKTICK_CLIENT_ID", ""),
			ClientSecret: getEnv("TICKTICK_CLIENT_SECRET", ""),
			AccessToken:  getEnv("TICKTICK_ACCESS_TOKEN", ""),
			RefreshToken: getEnv("TICKTICK_REFRESH_TOKEN", ""),
			// 其他配置使用默认值，不从环境变量读取
			BaseURL:     "https://api.dida365.com/open/v1",
			TokenURL:    "https://dida365.com/oauth/token",
			AuthURL:     "https://dida365.com/oauth/authorize",
			RedirectURL: "http://localhost:8000/callback",
			Timeout:     30 * time.Second,
		},
		Server: ServerConfig{
			Name:    "TickTick MCP Server",
			Version: "1.0.0",
			Port:    8000,
		},
		Log: LogConfig{
			Level:    "info",
			FilePath: "log.txt",
		},
	}

	// 验证必要的配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证 OAuth2 认证必需的配置
	if c.TickTick.ClientID == "" {
		return errors.New(errors.ErrInvalidCredentials, "TICKTICK_CLIENT_ID is required")
	}

	if c.TickTick.ClientSecret == "" {
		return errors.New(errors.ErrInvalidCredentials, "TICKTICK_CLIENT_SECRET is required")
	}

	// 验证 API 端点配置（这些有默认值，但仍需验证）
	if c.TickTick.BaseURL == "" {
		return errors.New(errors.ErrConfigLoad, "TICKTICK_BASE_URL is required")
	}

	if c.TickTick.TokenURL == "" {
		return errors.New(errors.ErrConfigLoad, "TICKTICK_TOKEN_URL is required")
	}

	if c.TickTick.AuthURL == "" {
		return errors.New(errors.ErrConfigLoad, "TICKTICK_AUTH_URL is required")
	}

	// 注意：AccessToken 不再强制要求，因为可以通过 OAuth2 流程获取
	// 如果没有 AccessToken，应用程序会引导用户完成 OAuth2 授权流程

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型的环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration 获取时间间隔类型的环境变量
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
