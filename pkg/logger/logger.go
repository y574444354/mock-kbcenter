package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zgsm/mock-kbcenter/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	asynqLogger asynq.Logger
)

// InitLogger 初始化日志系统
func InitLogger(cfg config.Log) error {
	// 确保日志目录存在
	if cfg.OutputPath != "" {
		dir := filepath.Dir(cfg.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	if cfg.ErrorPath != "" {
		dir := filepath.Dir(cfg.ErrorPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置主日志轮转
	var writeSyncer zapcore.WriteSyncer
	if cfg.OutputPath != "" {
		// 替换日志文件名中的日期占位符
		// Go语言使用"2006-01-02"作为日期格式模板，代表YYYY-MM-DD格式
		outputPath := strings.ReplaceAll(cfg.OutputPath, "{yyyy-mm-dd}", time.Now().Format("2006-01-02"))
		lumberJackLogger := &lumberjack.Logger{
			Filename:   outputPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true, // 使用本地时间
		}
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(lumberJackLogger),
			zapcore.AddSync(os.Stdout), // 保持控制台日志输出
		)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 配置错误日志轮转
	var errorWriteSyncer zapcore.WriteSyncer
	if cfg.ErrorPath != "" {
		// 替换错误日志文件名中的日期占位符
		// Go语言使用"2006-01-02"作为日期格式模板，代表YYYY-MM-DD格式
		errorPath := strings.ReplaceAll(cfg.ErrorPath, "{yyyy-mm-dd}", time.Now().Format("2006-01-02"))
		errorLogger := &lumberjack.Logger{
			Filename:   errorPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true, // 使用本地时间
		}
		errorWriteSyncer = zapcore.AddSync(errorLogger)
	}

	// 创建主日志核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建错误日志核心
	var cores []zapcore.Core
	cores = append(cores, core)

	if errorWriteSyncer != nil {
		errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		errorCore := zapcore.NewCore(encoder, errorWriteSyncer, errorLevel)
		cores = append(cores, errorCore)
	}

	// 创建日志记录器
	logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))

	// 创建asynq日志适配器
	asynqLogger = &zapAsynqLogger{
		logger: logger,
	}

	return nil
}

// Debug 记录调试级别日志
func Debug(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Debugw(msg, keysAndValues...)
}

// Info 记录信息级别日志
func Info(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Infow(msg, keysAndValues...)
}

// Warn 记录警告级别日志
func Warn(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Warnw(msg, keysAndValues...)
}

// Error 记录错误级别日志
func Error(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Errorw(msg, keysAndValues...)
}

// DPanic 记录开发环境恐慌级别日志
func DPanic(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.DPanicw(msg, keysAndValues...)
}

// Panic 记录恐慌级别日志
func Panic(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Panicw(msg, keysAndValues...)
}

// Fatal 记录致命级别日志
func Fatal(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Fatalw(msg, keysAndValues...)
}

// Sync 同步日志
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// GetAsynqLogger 获取Asynq日志记录器
func GetAsynqLogger() asynq.Logger {
	return asynqLogger
}

// zapAsynqLogger 实现asynq.Logger接口
type zapAsynqLogger struct {
	logger *zap.Logger
}

func (l *zapAsynqLogger) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

func (l *zapAsynqLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *zapAsynqLogger) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *zapAsynqLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *zapAsynqLogger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}
