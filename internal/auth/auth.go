package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

type TickTickAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURl     string
	Scopes       []string
	Port         int
	Config       *oauth2.Config
}

// NewTickTickAuth 创建一个新的TickTick认证管理器
func NewTickTickAuth(clientID, clientSecret string) (*TickTickAuth, error) {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("clientID or clientSecret missing")
	}

	// 获取认证URL
	authURL := os.Getenv("TICKTICK_AUTH_URL")
	if authURL == "" {
		authURL = "https://dida365.com/oauth/authorize"
	}

	// 获取令牌URL
	tokenURL := os.Getenv("TICKTICK_TOKEN_URL")
	if tokenURL == "" {
		tokenURL = "https://dida365.com/oauth/token"
	}

	// 默认作用域
	scopes := []string{"tasks:read", "tasks:write"}
	// 创建OAuth2配置
	redirectURI := "http://localhost:8000/callback"
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	return &TickTickAuth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		AuthURL:      authURL,
		TokenURl:     tokenURL,
		Scopes:       scopes,
		Port:         8000,
		Config:       config,
	}, nil
}

func (a *TickTickAuth) StartAuthFlow() error {
	if a.ClientID == "" || a.ClientSecret == "" {
		return fmt.Errorf("client ID or client secret missing")
	}

	// 生成随机state 参数
	state := base64.StdEncoding.EncodeToString([]byte(time.Now().String()))
	authURL := a.Config.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "code"))
	fmt.Println(authURL)
	// 打开浏览器
	if err := open.Run(authURL); err != nil {
		log.Printf("Warning: Failed to open browser: %v", err)
		log.Println("Please open the URL manually.")
	}
	// 启动本地服务器处理回调

	code, err := a.startCallbackServer(state)
	if err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}

	// 交换授权码获取令牌
	codeStr, ok := code.(string)
	if !ok {
		return fmt.Errorf("code is not a string type")
	}
	return a.exchangeCodeForToken(codeStr)
}

func (a *TickTickAuth) startCallbackServer(expectedState string) (interface{}, interface{}) {
	var authCode string
	var serverErr error
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// 创建服务器
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", a.Port),
	}

	// 定义回调处理函数
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// 获取查询参数
		query := r.URL.Query()
		//receivedState := query.Get("state")
		code := query.Get("code")
		errorMsg := query.Get("error")

		// 检查错误
		if errorMsg != "" {
			errChan <- fmt.Errorf("authorization error: %s", errorMsg)
			http.Error(w, fmt.Sprintf("Authorization error: %s", errorMsg), http.StatusBadRequest)
			return
		}

		// 查看授权码
		if code == "" {
			errChan <- fmt.Errorf("missing authorization code")
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}
		// 成功获取授权码
		codeChan <- code

		// 返回成功页面
		w.Header().Add("Content-Type", "text/html")

		w.Write([]byte(`
		<html>
		<head>
			<title>TickTick MCP Server - Authentication Successful</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					max-width: 600px;
					margin: 0 auto;
					padding: 20px;
					text-align: center;
				}
				h1 {
					color: #4CAF50;
				}
				.box {
					border: 1px solid #ddd;
					border-radius: 5px;
					padding: 20px;
					margin-top: 20px;
					background-color: #f9f9f9;
				}
			</style>
		</head>
		<body>
			<h1>Authentication Successful!</h1>
			<div class="box">
				<p>You have successfully authenticated with TickTick.</p>
				<p>You can now close this window and return to the terminal.</p>
			</div>
		</body>
		</html>
		`))

		// 优雅关闭服务器
		go func() {
			time.Sleep(time.Second * 10)
			server.Shutdown(context.Background())
		}()
		//
	})

	// 启动服务器
	go func() {
		log.Printf("Waiting for authentication callback on port %d...", a.Port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// 等待授权码
	select {
	case authCode = <-codeChan:
		return authCode, nil
	case serverErr = <-errChan:
		return "", serverErr
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("authentication time out")

	}

}

func (a *TickTickAuth) exchangeCodeForToken(code string) error {
	// 使用授权码交换令牌
	token, err := a.Config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("token exchange failed: %v", err)
	}
	// 保存令牌到环境文件
	if err := a.saveTokensToEnv(token.AccessToken, token.RefreshToken); err != nil {
		return fmt.Errorf("error saving tokens: %w", err)
	}
	return nil
}

func (a *TickTickAuth) saveTokensToEnv(accessToken string, refreshToken string) error {
	//
	envMap, err := godotenv.Read()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	if envMap == nil {
		envMap = make(map[string]string)
	}

	//update
	envMap["TICKTICK_ACCESS_TOKEN"] = accessToken
	envMap["TICKTICK_REFRESH_TOKEN"] = refreshToken
	envMap["TICKTICK_CLIENT_ID"] = a.ClientID
	envMap["TICKTICK_CLIENT_SECRET"] = a.ClientSecret
	// 保存.env文件
	return godotenv.Write(envMap, ".env")
}
