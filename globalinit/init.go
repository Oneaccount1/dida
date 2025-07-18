package globalinit

import (
	"dida/internal/logger"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

var (
	Log     *logger.Logger
	EnvPath string = ".env" // 默认环境变量文件路径
)

// SetEnvPath 设置环境变量文件路径
func SetEnvPath(path string) {
	if path != "" {
		EnvPath = path
	}
}

func Init() error {
	if err := initLogger(); err != nil {
		fmt.Printf("日志启动失败! err: %v", err)
		return err
	}
	if err := initEnv(); err != nil {
		Log.Errorf("初始化环境变量失败: %v", err)
		return err
	}

	return nil
}

func initEnv() error {
	// 获取环境变量文件的绝对路径
	absEnvPath, err := filepath.Abs(EnvPath)
	if err != nil {
		return fmt.Errorf("获取环境变量文件绝对路径失败: %v", err)
	}

	if _, err := os.Stat(absEnvPath); err != nil {
		return fmt.Errorf("未找到环境变量文件 %s: %v", absEnvPath, err)
	}

	// 加载环境变量
	if err := godotenv.Load(absEnvPath); err != nil {
		return fmt.Errorf("加载环境变量文件失败 %s: %v", absEnvPath, err)
	}

	Log.Infof("成功加载环境变量文件: %s", absEnvPath)
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
