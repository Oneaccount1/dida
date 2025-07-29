package server

import (
	"dida/globalinit"
	"dida/internal/client"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
)

var ticktickClient *client.TickTickClient

// ensureClientInitialized 确保客户端已初始化，如果未初始化则尝试初始化
func ensureClientInitialized() error {
	if ticktickClient == nil {
		if !InitializeClient() {
			return fmt.Errorf("failed to initialize TickTick client. Please check your API credentials")
		}
	}
	return nil
}

func InitializeClient() bool {
	var err error
	logger := globalinit.GetLogger()

	ticktickClient, err = client.NewTickTickClient()
	if err != nil {
		logger.Errorf("Failed to initialize TickTick client: %v", err)
		logger.Error("Please check your .env file and ensure TICKTICK_CLIENT_ID and TICKTICK_CLIENT_SECRET are set.")
		return false
	}

	// 检查是否有访问令牌
	if ticktickClient.GetAccessToken() == "" {
		logger.Info("No access token found. Please use the oauth_authorize tool to complete OAuth2 authentication.")
		return true // 返回 true，因为客户端已成功初始化，只是需要授权
	}

	// 如果有访问令牌，测试 API 连接
	projects, err := ticktickClient.GetProjects()
	if err != nil {
		logger.Errorf("Failed to access TickTick API: %v", err)
		logger.Info("Your access token may have expired. Please use the oauth_authorize tool to refresh it.")
		return true // 返回 true，因为客户端已初始化，只是令牌可能过期
	}
	logger.Infof("Successfully connected to TickTick API with %d projects", len(projects))
	return true
}

func Start() error {
	logger := globalinit.GetLogger()
	// 初始化TickTick客户端
	if !InitializeClient() {
		return fmt.Errorf("fail to initialize TickTick client")
	}
	// 创建MCP服务器
	s := server.NewMCPServer(
		"TickTick MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	// 初始化所有Tools
	if err := InitAllTools(s); err != nil {
		logger.Errorf("Failed to initialize MCP tools: %v", err)
		return err
	}

	logger.Info("Starting TickTick MCP server...")

	// 启动服务器
	return server.ServeStdio(s)
}
