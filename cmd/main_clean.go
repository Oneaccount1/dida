package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"dida/cmd/app"
)

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	log.Println("Starting TickTick MCP Server (Clean Architecture)...")
	
	// 创建应用程序
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	
	// 创建一个channel来接收信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	
	// 启动服务器
	errChan := make(chan error, 1)
	go func() {
		err := application.Start()
		if err != nil {
			errChan <- err
		}
	}()
	
	// 等待信号或错误
	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigs:
		fmt.Printf("\nReceived signal: %v, shutting down...\n", sig)
	}
	
	log.Println("Server stopped")
}