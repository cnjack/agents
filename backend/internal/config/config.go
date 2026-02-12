package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	AI       AIConfig       `yaml:"ai"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// AIConfig AI 服务配置
type AIConfig struct {
	DefaultProvider string                    `yaml:"default_provider"`
	Providers       map[string]ProviderConfig `yaml:"providers"`
	Cache           CacheConfig               `yaml:"cache"`
}

// ProviderConfig AI 提供者配置
type ProviderConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Endpoint    string  `yaml:"endpoint"`
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
	Timeout     int     `yaml:"timeout"`
	MaxRetries  int     `yaml:"max_retries"`
	DelayMs     int     `yaml:"delay_ms"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled    bool `yaml:"enabled"`
	TTLSeconds int  `yaml:"ttl_seconds"`
	MaxEntries int  `yaml:"max_entries"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
	SQLite  SQLiteConfig  `yaml:"sqlite"`
}

// MongoDBConfig MongoDB 配置
type MongoDBConfig struct {
	Enabled  bool   `yaml:"enabled"`
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// SQLiteConfig SQLite 配置
type SQLiteConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgMu   sync.RWMutex
)

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	// 如果未指定路径，使用默认路径
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析 YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	setDefaults(&config)

	cfg = &config
	return cfg, nil
}

// Get 获取当前配置
func Get() *Config {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg
}

// getDefaultConfigPath 获取默认配置文件路径
func getDefaultConfigPath() string {
	// 按优先级查找配置文件
	searchPaths := []string{
		"./config/local.yaml",      // 本地配置（优先）
		"./config/config.yaml",     // 默认配置
		"./local.yaml",             // 根目录本地配置
		"./config.yaml",            // 根目录默认配置
	}

	// 获取可执行文件所在目录
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		searchPaths = append(searchPaths,
			filepath.Join(execDir, "config/local.yaml"),
			filepath.Join(execDir, "config/config.yaml"),
		)
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 默认返回 local.yaml 路径（即使不存在）
	return "./config/local.yaml"
}

// setDefaults 设置默认值
func setDefaults(c *Config) {
	// Server defaults
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}

	// AI defaults
	if c.AI.DefaultProvider == "" {
		c.AI.DefaultProvider = "mock"
	}

	// Provider defaults
	for name, provider := range c.AI.Providers {
		if provider.MaxTokens == 0 {
			provider.MaxTokens = 2048
		}
		if provider.Temperature == 0 {
			provider.Temperature = 0.7
		}
		if provider.Timeout == 0 {
			provider.Timeout = 60
		}
		if provider.MaxRetries == 0 {
			provider.MaxRetries = 3
		}
		c.AI.Providers[name] = provider
	}

	// Cache defaults
	if c.AI.Cache.TTLSeconds == 0 {
		c.AI.Cache.TTLSeconds = 3600
	}
	if c.AI.Cache.MaxEntries == 0 {
		c.AI.Cache.MaxEntries = 1000
	}

	// Logging defaults
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}
	if c.Logging.Output == "" {
		c.Logging.Output = "stdout"
	}
}

// GetProviderConfig 获取指定提供者的配置
func (c *Config) GetProviderConfig(name string) (ProviderConfig, bool) {
	config, ok := c.AI.Providers[name]
	return config, ok
}

// IsProviderEnabled 检查提供者是否启用
func (c *Config) IsProviderEnabled(name string) bool {
	config, ok := c.AI.Providers[name]
	return ok && config.Enabled
}
