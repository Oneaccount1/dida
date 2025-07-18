package main

import (
	"bufio"
	"dida/globalinit"
	"dida/internal/auth"
	"dida/internal/server"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	envPath = flag.String("env", ".env", "指定.env文件的路径")
)

func main() {
	flag.Parse()

	// 检查是否是认证命令
	if len(os.Args) > 1 && os.Args[1] == "auth" {
		authPath := ".env"
		// 检查是否有额外的参数指定env路径
		if len(os.Args) > 2 {
			authPath = os.Args[2]
		}
		authFlow(authPath)
		return
	}

	// 获取环境变量文件的绝对路径
	absEnvPath, err := filepath.Abs(*envPath)
	if err != nil {
		fmt.Printf("获取环境变量文件绝对路径失败: %v\n", err)
		time.Sleep(3 * time.Second)
		return
	}

	// 将环境变量文件路径设置到环境变量中，供client使用
	os.Setenv("TICKTICK_ENV_PATH", absEnvPath)

	// 设置环境变量文件路径
	globalinit.SetEnvPath(*envPath)

	// Otherwise start the server
	// 初始化
	if err := globalinit.Init(); err != nil {
		fmt.Printf("初始化失败:%v程序退出3s后退出...", err)
		time.Sleep(3 * time.Second)
		return
	}

	Log := globalinit.GetLogger()
	if Log == nil {
		return
	}

	Log.Info("Starting TickTick MCP Server")
	Log.Infof("使用环境变量文件: %s", absEnvPath)

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

// 加载指定路径的环境变量文件
func loadEnvFile(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("获取环境变量文件绝对路径失败: %v\n", err)
		return
	}

	// 将环境变量文件路径设置到环境变量中，供client使用
	os.Setenv("TICKTICK_ENV_PATH", absPath)

	if err := godotenv.Load(absPath); err != nil {
		fmt.Printf("加载环境变量文件失败 %s: %v\n", absPath, err)
	} else {
		fmt.Printf("成功加载环境变量文件: %s\n", absPath)
	}
}

func authFlow(envPath string) {
	// 认证模式下，直接输出到控制台是可以的，因为不会干扰MCP协议
	loadEnvFile(envPath)

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

	// 获取认证信息
	clientID = tickAuth.GetClientID()
	clientSecret = tickAuth.GetClientSecret()

	// 保存凭证到指定的.env文件
	accessToken := os.Getenv("TICKTICK_ACCESS_TOKEN")
	refreshToken := os.Getenv("TICKTICK_REFRESH_TOKEN")

	// 使用auth包的方法保存令牌到指定路径
	err = tickAuth.SaveTokensToCustomEnv(accessToken, refreshToken, envPath)
	if err != nil {
		fmt.Printf("保存凭证失败: %v\n", err)
		return
	}

	fmt.Println("认证成功!, 你现在可以启动MCP服务器！")
	fmt.Printf("凭证已保存到: %s\n", envPath)
}

func Load() (string, string) {
	// 认证模式下，直接输出到控制台是可以的，因为不会干扰MCP协议
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
