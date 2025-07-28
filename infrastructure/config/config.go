package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
	
	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	// TickTick API配置
	TickTick TickTickConfig
	
	// 认证配置
	Auth AuthConfig
	
	// 服务器配置
	Server ServerConfig
	
	// 存储配置
	Storage StorageConfig
	
	// 日志配置
	Log LogConfig
}

// TickTickConfig TickTick API配置
type TickTickConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
	Scopes       []string
}

// AuthConfig 认证配置
type AuthConfig struct {
	TokenFilePath string
	Port          int
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name    string
	Version string
	Timeout time.Duration
}

// StorageConfig 存储配置
type StorageConfig struct {
	DataDir string
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string
	Format string
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 加载环境变量文件
	if err := godotenv.Load(); err != nil {
		// 忽略.env文件不存在的错误
	}
	
	config := &Config{
		TickTick: TickTickConfig{
			BaseURL:      getEnv("TICKTICK_BASE_URL", "https://api.dida365.com/open/v1"),
			ClientID:     getEnv("TICKTICK_CLIENT_ID", ""),
			ClientSecret: getEnv("TICKTICK_CLIENT_SECRET", ""),
			AuthURL:      getEnv("TICKTICK_AUTH_URL", "https://dida365.com/oauth/authorize"),
			TokenURL:     getEnv("TICKTICK_TOKEN_URL", "https://dida365.com/oauth/token"),
			RedirectURL:  getEnv("TICKTICK_REDIRECT_URL", "http://localhost:8080/callback"),
			Scopes:       []string{"tasks:read", "tasks:write"},
		},
		Auth: AuthConfig{
			TokenFilePath: getEnv("TICKTICK_TOKEN_FILE", filepath.Join(os.TempDir(), "ticktick_auth.json")),
			Port:          getEnvAsInt("TICKTICK_AUTH_PORT", 8080),
		},
		Server: ServerConfig{
			Name:    "TickTick MCP Server",
			Version: "1.0.0",
			Timeout: time.Duration(getEnvAsInt("SERVER_TIMEOUT", 30)) * time.Second,
		},
		Storage: StorageConfig{
			DataDir: getEnv("DATA_DIR", filepath.Join(os.TempDir(), "ticktick-mcp")),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
	
	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.TickTick.ClientID == "" {
		return &ConfigError{Field: "TICKTICK_CLIENT_ID", Message: "missing required field"}
	}
	if c.TickTick.ClientSecret == "" {
		return &ConfigError{Field: "TICKTICK_CLIENT_SECRET", Message: "missing required field"}
	}
	
	return nil
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Field + " - " + e.Message
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}