package main

import (
	"dida/internal/logger"
	"dida/internal/server"
	"fmt"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

var Log *logger.Logger

func main() {
	// 初始化日志
	zapLog, err := logger.NewLogger("Log.txt", zapcore.InfoLevel)
	if err != nil {
		fmt.Printf("日志启动失败! err: %v", err)
		return
	}
	Log = zapLog

	Log.Info("Starting TickTick MCP Server")

	// 创建一个channel来接收信号
	sigs := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// 启动服务器
	go func() {
		err := server.Start()
		if err != nil {
			errChan <- err
		}
	}()

	// 等待信号或错误
	select {
	case err := <-errChan:
		Log.Errorf("Server error: %v", err)
	case sig := <-sigs:
		Log.Infof("Received signal: %v shutting down...", sig)
	}
}
