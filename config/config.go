package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"gopkg.in/yaml.v3"
)

// Log configuration
type Log struct {
	Level      string `yaml:"level"`  // debug, info, warn, error, dpanic, panic, fatal
	Format     string `yaml:"format"` // json, console
	OutputPath string `yaml:"output_path"`
	ErrorPath  string `yaml:"error_path"` // Path for storing error logs separately
	MaxSize    int    `yaml:"max_size"`   // MB
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"` // days
	Compress   bool   `yaml:"compress"`
}

type Database struct {
	Type     string `yaml:"type"` // mysql, postgres, sqlite
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"timezone"`
	Enabled  bool   `yaml:"enabled"` // Whether to enable database
}

// Config application configuration structure
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"` // debug, release, test
	} `yaml:"server"`

	Database Database `yaml:"database"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		Enabled  bool   `yaml:"enabled"` // Whether to enable Redis
	} `yaml:"redis"`

	// Asynchronous task queue configuration
	Asynq struct {
		Concurrency   int            `yaml:"concurrency"`     // Number of concurrent workers
		RetryCount    int            `yaml:"retry_count"`     // Maximum retry count
		RetryDelay    int            `yaml:"retry_delay"`     // Retry delay in seconds
		RedisPoolSize int            `yaml:"redis_pool_size"` // Redis connection pool size
		Queues        map[string]int `yaml:"queues"`          // Queue priorities
		Log           Log            `yaml:"log"`             // Asynq log configuration
		Enabled       bool           `yaml:"enabled"`         // Whether to enable Asynq
	} `yaml:"asynq"`

	Log Log `yaml:"log"`

	I18n struct {
		DefaultLocale string `yaml:"default_locale"`
		BundlePath    string `yaml:"bundle_path"`
	} `yaml:"i18n"`

	// HTTPClient HTTP client configuration
	HTTPClient struct {
		// Default timeout in seconds
		Timeout int `yaml:"timeout"`
		// Default retry count
		MaxRetries int `yaml:"max_retries"`
		// Retry delay in seconds
		RetryDelay int `yaml:"retry_delay"`
		// Whether to enable request logging
		EnableRequestLog bool `yaml:"enable_request_log"`
		// Whether to enable response logging
		EnableResponseLog bool `yaml:"enable_response_log"`
		// Default request headers
		Headers map[string]string `yaml:"headers"`
		// Proxy URL
		ProxyURL string `yaml:"proxy_url"`
		// TLS configuration
		InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
		// Dependent service configurations
		Services map[string]ServiceConfig `yaml:"services"`
	} `yaml:"http_client"`
}

// ServiceConfig dependent service configuration
type ServiceConfig struct {
	// Service base URL
	BaseURL string `yaml:"base_url"`
	// Timeout in seconds, overrides default
	Timeout int `yaml:"timeout"`
	// Retry count, overrides default
	MaxRetries int `yaml:"max_retries"`
	// Authentication type: none, basic, bearer, custom
	AuthType string `yaml:"auth_type"`
	// Basic auth username
	Username string `yaml:"username"`
	// Basic auth password
	Password string `yaml:"password"`
	// Bearer auth token
	Token string `yaml:"token"`
	// Custom auth header
	AuthHeader string `yaml:"auth_header"`
	// Request headers, merged with default headers
	Headers map[string]string `yaml:"headers"`
	// Valid status codes, empty means default (2xx)
	ValidStatusCodes []int `yaml:"valid_status_codes"`
}

//go:embed config.yaml
var defaultConfig embed.FS

var (
	config *Config
	once   sync.Once
)

// mergeConfig recursively merges configurations, target is destination config, source is source config
func mergeConfig(target, source interface{}) {
	targetVal := reflect.ValueOf(target).Elem()
	sourceVal := reflect.ValueOf(source).Elem()

	for i := 0; i < sourceVal.NumField(); i++ {
		sourceField := sourceVal.Field(i)
		if !sourceField.IsZero() {
			targetField := targetVal.Field(i)
			switch targetField.Kind() {
			case reflect.Struct:
				// Recursively process nested structs
				mergeConfig(targetField.Addr().Interface(), sourceField.Addr().Interface())
			case reflect.Map:
				// Merge map type fields
				if targetField.IsNil() {
					targetField.Set(reflect.MakeMap(targetField.Type()))
				}
				for _, key := range sourceField.MapKeys() {
					targetField.SetMapIndex(key, sourceField.MapIndex(key))
				}
			default:
				// Directly assign basic type fields
				targetField.Set(sourceField)
			}
		}
	}
}

// LoadConfigWithDefault loads default config file and tries to override with local config
func LoadConfigWithDefault() (err error) {
	once.Do(func() {
		config = &Config{}

		// Get base path
		var basePath string
		basePath, err = os.Getwd()
		if err != nil {
			basePath = filepath.Dir(os.Args[0])
		}
		fmt.Println("Base path:", basePath)

		// Read embedded default config file
		var defaultData []byte
		defaultData, err = defaultConfig.ReadFile("config.yaml")
		if err != nil {
			// If embedded file read fails, try reading from filesystem
			defaultPath := filepath.Join(basePath, "config", "config.yaml")
			defaultData, err = os.ReadFile(defaultPath)
			if err != nil {
				return
			}
		}

		// Parse default config
		if err = yaml.Unmarshal(defaultData, config); err != nil {
			return
		}

		// Try to read local config file for override
		localPath := filepath.Join(basePath, "config", "config.local.yaml")
		localData, readLocalErr := os.ReadFile(localPath)
		if readLocalErr == nil {
			// Local config file exists, parse and merge with default config
			fmt.Println("Using local config file:", localPath)
			localConfig := &Config{}
			if unmarshalErr := yaml.Unmarshal(localData, localConfig); unmarshalErr == nil {
				mergeConfig(config, localConfig)
			}
		}
	})
	return err
}

// GetConfig gets configuration instance
func GetConfig() *Config {
	if config == nil {
		// If config is not loaded, use default config
		once.Do(func() {
			config = &Config{}
			// Set default values
			config.Server.Port = 8080
			config.Server.Mode = "release"
			config.Database.Enabled = true
			config.Redis.Enabled = true
			config.Asynq.Enabled = true
		})
	}
	return config
}
