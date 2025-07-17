package main

import (
	"dida/internal/auth"
	"dida/internal/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	clientID := "4Fj1OR96WXPQPlu1aK"
	clientSecret := "7kZL%&l7&zi*22w5Fk+NEk(F77_rDB1o"
	// 创建认证管理器
	ticktickAuth := auth.NewTickTickAuth(clientID, clientSecret)

	// 启动认证流程
	result, err := ticktickAuth.StartAuthFlow()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	log.Println(result)
	log.Println("You can now start the MCP server.")

	// 以下是原有的服务器启动流程
	log.Println("Starting TickTick MCP Server...")

	// 创建一个channel来接收信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	errChan := make(chan error, 1)
	go func() {
		err := server.Start()
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
