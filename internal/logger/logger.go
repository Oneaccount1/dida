package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 封装zap.Logger
type Logger struct {
	zapLogger *zap.SugaredLogger
}

// NewLogger 创建一个新的日志记录器，带有文件输出
func NewLogger(logFile string, level zapcore.Level) (*Logger, error) {
	// 确保日志目录存在
	// Dir返回路径的目录部分
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 打开日志文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	// 创建日志配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建core
	core := zapcore.NewTee(
		// 控制台输出
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
		// 文件输出
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			level,
		),
	)

	// 创建Logger
	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{zapLogger: zapLogger.Sugar()}, nil
}

// Debug 打印调试级别日志
func (l *Logger) Debug(args ...interface{}) {
	l.zapLogger.Debug(args...)
}

// Debugf 打印格式化调试级别日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zapLogger.Debugf(format, args...)
}

// Info 打印信息级别日志
func (l *Logger) Info(args ...interface{}) {
	l.zapLogger.Info(args...)
}

// Infof 打印格式化信息级别日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.zapLogger.Infof(format, args...)
}

// Warn 打印警告级别日志
func (l *Logger) Warn(args ...interface{}) {
	l.zapLogger.Warn(args...)
}

// Warnf 打印格式化警告级别日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zapLogger.Warnf(format, args...)
}

// Error 打印错误级别日志
func (l *Logger) Error(args ...interface{}) {
	l.zapLogger.Error(args...)
}

// Errorf 打印格式化错误级别日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zapLogger.Errorf(format, args...)
}

// Fatal 打印致命级别日志，并退出程序
func (l *Logger) Fatal(args ...interface{}) {
	l.zapLogger.Fatal(args...)
}

// Fatalf 打印格式化致命级别日志，并退出程序
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.zapLogger.Fatalf(format, args...)
}
