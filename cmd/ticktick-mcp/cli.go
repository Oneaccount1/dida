package main

import (
	"bufio"
	"dida/internal/auth"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func main() {

	// 加载
	if err := godotenv.Load(); err != nil {
		fmt.Printf("加载环境变量有误: %v\n", err)
	}
	clientID := os.Getenv("TICKTICK_CLIENT_ID")
	clientSecret := os.Getenv("TICKTICK_CLIENT_SECRET")

	if clientID != "" && clientSecret != "" {
		fmt.Println("凭证已经存在， 无需重复验证")
		fmt.Printf("验证信息:\nclientID: %s \nclientSecret: %s \n", clientID, clientSecret)
		return
	}

	tickAuth, err := auth.NewTickTickAuth(Load())

	if err != nil {
		fmt.Printf("初始化认证失败: %v", err)
		return
	}

	// 开始认证
	err = tickAuth.StartAuthFlow()

	if err != nil {
		fmt.Printf("认证失败: %v", err)
		return
	}
	fmt.Println("认证成功!, 你现在可以启动MCP服务器！")
}

func Load() (string, string) {

	fmt.Println(`
	╔════════════════════════════════════════════════╗
	║       TickTick MCP Server Authentication       ║
	╚════════════════════════════════════════════════╝
		
	This utility will help you authenticate with TickTick
	and obtain the necessary access tokens for the TickTick MCP server.
	
	Before you begin, you will need:
	1. A TickTick account (https://ticktick.com)
	2. A registered TickTick API application (https://developer.ticktick.com)
	3. Your Client ID and Client Secret from the TickTick Developer Center`)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请输入clientID: ")
	clientID, _ := reader.ReadString('\n')
	clientID = strings.TrimSpace(clientID)

	fmt.Println("请输入clientSecret: ")
	clientSecret, _ := reader.ReadString('\n')
	clientSecret = strings.TrimSpace(clientSecret)
	fmt.Printf("验证信息:\nclientID: %s \nclientSecret: %s \n", clientID, clientSecret)
	return clientID, clientSecret
}
