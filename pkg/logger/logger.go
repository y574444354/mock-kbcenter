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

// InitLogger initialize logging system
func InitLogger(cfg config.Log) error {
	// Ensure log directory exists
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

	// Set log level
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

	// Configure encoder
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

	// Configure output
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure main log rotation
	var writeSyncer zapcore.WriteSyncer
	if cfg.OutputPath != "" {
		// Replace date placeholder in log filename
		// Go uses "2006-01-02" as date format template, representing YYYY-MM-DD format
		outputPath := strings.ReplaceAll(cfg.OutputPath, "{yyyy-mm-dd}", time.Now().Format("2006-01-02"))
		lumberJackLogger := &lumberjack.Logger{
			Filename:   outputPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true, // Use local time
		}
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(lumberJackLogger),
			zapcore.AddSync(os.Stdout), // Keep console log output
		)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Configure error log rotation
	var errorWriteSyncer zapcore.WriteSyncer
	if cfg.ErrorPath != "" {
		// Replace date placeholder in error log filename
		// Go uses "2006-01-02" as date format template, representing YYYY-MM-DD format
		errorPath := strings.ReplaceAll(cfg.ErrorPath, "{yyyy-mm-dd}", time.Now().Format("2006-01-02"))
		errorLogger := &lumberjack.Logger{
			Filename:   errorPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true, // Use local time
		}
		errorWriteSyncer = zapcore.AddSync(errorLogger)
	}

	// Create main log core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create error log core
	var cores []zapcore.Core
	cores = append(cores, core)

	if errorWriteSyncer != nil {
		errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		errorCore := zapcore.NewCore(encoder, errorWriteSyncer, errorLevel)
		cores = append(cores, errorCore)
	}

	// Create logger
	logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))

	// Create asynq logger adapter
	asynqLogger = &zapAsynqLogger{
		logger: logger,
	}

	return nil
}

// Debug log debug level message
func Debug(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Debugw(msg, keysAndValues...)
}

// Info log info level message
func Info(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Infow(msg, keysAndValues...)
}

// Warn log warning level message
func Warn(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Warnw(msg, keysAndValues...)
}

// Error log error level message
func Error(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Errorw(msg, keysAndValues...)
}

// DPanic log development panic level message
func DPanic(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.DPanicw(msg, keysAndValues...)
}

// Panic log panic level message
func Panic(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Panicw(msg, keysAndValues...)
}

// Fatal log fatal level message
func Fatal(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}
	sugar := logger.Sugar()
	sugar.Fatalw(msg, keysAndValues...)
}

// Sync flush log buffer
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// GetAsynqLogger get asynq logger
func GetAsynqLogger() asynq.Logger {
	return asynqLogger
}

// zapAsynqLogger implements asynq.Logger interface
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
