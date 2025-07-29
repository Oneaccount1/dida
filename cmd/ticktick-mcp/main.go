package main

import (
	"dida/globalinit"
	"dida/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

// initializeEnvironment 初始化环境变量和配置
func initializeEnvironment() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 初始化全局组件
	if err := globalinit.Init(); err != nil {
		return err
	}

	return nil
}

func main() {
	// 初始化环境变量和配置
	if err := initializeEnvironment(); err != nil {
		log.Printf("初始化失败: %v", err)
		return
	}

	// 获取日志器
	logger := globalinit.GetLogger()
	if logger == nil {
		log.Println("日志器初始化失败")
		return
	}

	logger.Info("TickTick MCP Server initialized successfully")

	// 创建信号处理通道
	sigs := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		logger.Info("Starting MCP server...")
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// 等待信号或错误
	select {
	case err := <-errChan:
		logger.Errorf("Server error: %v", err)
	case sig := <-sigs:
		logger.Infof("Received signal: %v, shutting down...", sig)
	}

	logger.Info("TickTick MCP Server stopped")
}
