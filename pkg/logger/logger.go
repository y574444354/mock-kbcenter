package logger

import (
	"os"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/zgsm/review-manager/config"
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

	// 配置日志轮转
	var writeSyncer zapcore.WriteSyncer
	if cfg.OutputPath != "" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(lumberJackLogger),
			zapcore.AddSync(os.Stdout),
		)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建日志记录器
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

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
