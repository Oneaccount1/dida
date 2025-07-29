package globalinit

import (
	"dida/internal/logger"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

var (
	Log *logger.Logger
)

func Init() error {
	if err := initLogger(); err != nil {
		fmt.Printf("日志启动失败! err: %v", err)
		return err
	}
	return nil
}

func initLogger() error {
	execDir, err := os.Getwd()
	if err != nil {
		return err
	}

	logPath := filepath.Join(execDir, "log.txt")

	zapLog, err := logger.NewLogger(logPath, zapcore.InfoLevel)
	if err != nil {
		return err
	}
	Log = zapLog
	return nil
}

func GetLogger() *logger.Logger {
	return Log
}
