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

// Log 日志配置
type Log struct {
	Level      string `yaml:"level"`  // debug, info, warn, error, dpanic, panic, fatal
	Format     string `yaml:"format"` // json, console
	OutputPath string `yaml:"output_path"`
	ErrorPath  string `yaml:"error_path"` // 错误日志单独存储路径
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
	Enabled  bool   `yaml:"enabled"` // 是否启用数据库
}

// Config 应用配置结构
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
		Enabled  bool   `yaml:"enabled"` // 是否启用Redis
	} `yaml:"redis"`

	// 异步任务队列配置
	Asynq struct {
		Concurrency   int            `yaml:"concurrency"`     // 并发worker数量
		RetryCount    int            `yaml:"retry_count"`     // 最大重试次数
		RetryDelay    int            `yaml:"retry_delay"`     // 重试延迟(秒)
		RedisPoolSize int            `yaml:"redis_pool_size"` // Redis连接池大小
		Queues        map[string]int `yaml:"queues"`          // 队列优先级
		Log           Log            `yaml:"log"`             // Asynq日志配置
		Enabled       bool           `yaml:"enabled"`         // 是否启用Asynq
	} `yaml:"asynq"`

	Log Log `yaml:"log"`

	I18n struct {
		DefaultLocale string `yaml:"default_locale"`
		BundlePath    string `yaml:"bundle_path"`
	} `yaml:"i18n"`

	// HTTPClient HTTP客户端配置
	HTTPClient struct {
		// 默认超时设置（秒）
		Timeout int `yaml:"timeout"`
		// 默认重试次数
		MaxRetries int `yaml:"max_retries"`
		// 重试延迟（秒）
		RetryDelay int `yaml:"retry_delay"`
		// 是否启用请求日志
		EnableRequestLog bool `yaml:"enable_request_log"`
		// 是否启用响应日志
		EnableResponseLog bool `yaml:"enable_response_log"`
		// 默认请求头
		Headers map[string]string `yaml:"headers"`
		// 代理URL
		ProxyURL string `yaml:"proxy_url"`
		// TLS配置
		InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
		// 依赖服务配置
		Services map[string]ServiceConfig `yaml:"services"`
	} `yaml:"http_client"`
}

// ServiceConfig 依赖服务配置
type ServiceConfig struct {
	// 服务基础URL
	BaseURL string `yaml:"base_url"`
	// 超时设置（秒），覆盖默认值
	Timeout int `yaml:"timeout"`
	// 重试次数，覆盖默认值
	MaxRetries int `yaml:"max_retries"`
	// 鉴权类型：none, basic, bearer, custom
	AuthType string `yaml:"auth_type"`
	// Basic认证用户名
	Username string `yaml:"username"`
	// Basic认证密码
	Password string `yaml:"password"`
	// Bearer认证Token
	Token string `yaml:"token"`
	// 自定义认证头
	AuthHeader string `yaml:"auth_header"`
	// 请求头，与默认请求头合并
	Headers map[string]string `yaml:"headers"`
	// 有效的状态码列表，为空则使用默认规则(2xx)
	ValidStatusCodes []int `yaml:"valid_status_codes"`
}

//go:embed config.yaml
var defaultConfig embed.FS

var (
	config *Config
	once   sync.Once
)

// mergeConfig 递归合并配置，target 是目标配置，source 是源配置
func mergeConfig(target, source interface{}) {
	targetVal := reflect.ValueOf(target).Elem()
	sourceVal := reflect.ValueOf(source).Elem()

	for i := 0; i < sourceVal.NumField(); i++ {
		sourceField := sourceVal.Field(i)
		if !sourceField.IsZero() {
			targetField := targetVal.Field(i)
			switch targetField.Kind() {
			case reflect.Struct:
				// 递归处理嵌套结构体
				mergeConfig(targetField.Addr().Interface(), sourceField.Addr().Interface())
			case reflect.Map:
				// 合并 map 类型字段
				if targetField.IsNil() {
					targetField.Set(reflect.MakeMap(targetField.Type()))
				}
				for _, key := range sourceField.MapKeys() {
					targetField.SetMapIndex(key, sourceField.MapIndex(key))
				}
			default:
				// 直接赋值基本类型字段
				targetField.Set(sourceField)
			}
		}
	}
}

// LoadConfigWithDefault 加载默认配置文件，并尝试用本地配置覆盖
func LoadConfigWithDefault() error {
	var err error
	once.Do(func() {
		config = &Config{}

		// 获取基础路径
		basePath, err := os.Getwd()
		if err != nil {
			basePath = filepath.Dir(os.Args[0])
		}
		fmt.Println("Base path:", basePath)

		// 读取嵌入的默认配置文件
		defaultData, readErr := defaultConfig.ReadFile("config.yaml")
		if readErr != nil {
			// 如果嵌入文件读取失败，尝试从文件系统读取
			defaultPath := filepath.Join(basePath, "config", "config.yaml")
			defaultData, readErr = os.ReadFile(defaultPath)
			if readErr != nil {
				err = readErr
				return
			}
		}

		// 解析默认配置
		if unmarshalErr := yaml.Unmarshal(defaultData, config); unmarshalErr != nil {
			err = unmarshalErr
			return
		}

		// 尝试读取本地配置文件进行覆盖
		localPath := filepath.Join(basePath, "config", "config.local.yaml")
		localData, readLocalErr := os.ReadFile(localPath)
		if readLocalErr == nil {
			// 本地配置文件存在，解析并合并到默认配置
			fmt.Println("Using local config file:", localPath)
			localConfig := &Config{}
			if unmarshalErr := yaml.Unmarshal(localData, localConfig); unmarshalErr == nil {
				mergeConfig(config, localConfig)
			}
		}

		// 环境变量覆盖
		// if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		// 	// 可以添加环境变量覆盖逻辑
		// }
	})
	return err
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	if config == nil {
		// 如果配置未加载，使用默认配置
		config = &Config{}
		// 设置默认值
		config.Server.Port = 8080
		config.Server.Mode = "debug"
		config.Database.Enabled = true
		config.Redis.Enabled = true
		config.Asynq.Enabled = true
	}
	return config
}
