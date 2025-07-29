package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"dida/internal/errors"
	"github.com/joho/godotenv"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

const (
	defaultLocation    = ".env"
	defaultAuthURL     = "https://dida365.com/oauth/authorize"
	defaultTokenURL    = "https://dida365.com/oauth/token"
	defaultRedirectURI = "http://localhost:8000/callback"
)

// TokenResponse 令牌响应结构
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TickTickAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	Scopes       []string
	Port         int
	Config       *oauth2.Config
	EnvPath      string
}

// NewTickTickAuth 创建一个新的TickTick认证管理器
func NewTickTickAuth(clientID, clientSecret string) (*TickTickAuth, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("clientID or clientSecret missing")
	}

	// 获取认证URL
	authURL := os.Getenv("TICKTICK_AUTH_URL")
	if authURL == "" {
		authURL = defaultAuthURL
	}

	// 获取令牌URL
	tokenURL := os.Getenv("TICKTICK_TOKEN_URL")
	if tokenURL == "" {
		tokenURL = defaultTokenURL
	}

	// 默认作用域
	scopes := []string{"tasks:read", "tasks:write"}
	// 创建OAuth2配置
	redirectURI := defaultRedirectURI
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
		TokenURL:     tokenURL,
		Scopes:       scopes,
		Port:         8000,
		Config:       config,
	}, nil
}

// GetAuthURL 生成OAuth2授权URL
func (a *TickTickAuth) GetAuthURL() string {
	// 生成随机state参数
	state := base64.StdEncoding.EncodeToString([]byte(time.Now().String()))
	return a.Config.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "code"))
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
		fmt.Printf("Warning: Failed to open browser: %v\n", err)
		fmt.Println("Please open the URL manually.")
	}
	// 启动本地服务器处理回调

	code, err := a.startCallbackServer()
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

func (a *TickTickAuth) startCallbackServer() (interface{}, error) {
	var authCode string
	var serverErr error
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// 创建独立的多路复用器，避免路由冲突
	mux := http.NewServeMux()

	// 创建服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: mux,
	}

	// 定义回调处理函数
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
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

		// 检查授权码
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
		fmt.Printf("Waiting for authentication callback on port %d...\n", a.Port)
		if err := server.ListenAndServe(); !stderrors.Is(err, http.ErrServerClosed) {
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
		return "", fmt.Errorf("authentication timeout")

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

func (a *TickTickAuth) GetClientID() string {
	return a.ClientID
}

func (a *TickTickAuth) GetClientSecret() string {
	return a.ClientSecret
}

func (a *TickTickAuth) saveTokensToEnv(accessToken string, refreshToken string) error {
	// 直接使用 .env 文件
	envPath := defaultLocation

	// 读取现有的环境变量
	envMap, err := godotenv.Read(envPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	if envMap == nil {
		envMap = make(map[string]string)
	}

	// 更新令牌
	envMap["TICKTICK_ACCESS_TOKEN"] = accessToken
	envMap["TICKTICK_REFRESH_TOKEN"] = refreshToken
	envMap["TICKTICK_CLIENT_ID"] = a.ClientID
	envMap["TICKTICK_CLIENT_SECRET"] = a.ClientSecret

	// 保存.env文件
	return godotenv.Write(envMap, envPath)
}

// RefreshAccessToken 刷新访问令牌
func (a *TickTickAuth) RefreshAccessToken(currentRefreshToken string) (*TokenResponse, error) {
	if currentRefreshToken == "" {
		return nil, errors.New(errors.ErrTokenRefreshFailed, "no refresh token available")
	}
	if a.ClientID == "" || a.ClientSecret == "" {
		return nil, errors.New(errors.ErrInvalidCredentials, "client ID or client secret missing")
	}

	// 准备token请求
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", currentRefreshToken)

	// 准备Basic Auth验证
	authStr := fmt.Sprintf("%s:%s", a.ClientID, a.ClientSecret)
	authBytes := []byte(authStr)
	authB64 := base64.StdEncoding.EncodeToString(authBytes)

	// 创建请求
	request, err := http.NewRequest("POST", a.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrapf(errors.ErrAPIRequest, err, "failed to create token refresh request")
	}
	request.Header.Add("Authorization", "Basic "+authB64)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrAPIRequest, err, "failed to send token refresh request")
	}
	defer response.Body.Close()

	// 检查响应状态
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, errors.Newf(errors.ErrAPIResponse, "token refresh failed with status %s: %s", response.Status, string(body))
	}

	// 解析响应
	var tokens TokenResponse
	if err = json.NewDecoder(response.Body).Decode(&tokens); err != nil {
		return nil, errors.Wrapf(errors.ErrDataUnmarshal, err, "failed to parse token refresh response")
	}

	// 保存新令牌到环境文件
	refreshTokenToSave := tokens.RefreshToken
	if refreshTokenToSave == "" {
		// 如果响应中没有新的refresh token，使用原来的
		refreshTokenToSave = currentRefreshToken
	}

	if err := a.saveTokensToEnv(tokens.AccessToken, refreshTokenToSave); err != nil {
		return nil, errors.Wrapf(errors.ErrConfigLoad, err, "failed to save refreshed tokens")
	}

	return &tokens, nil
}
