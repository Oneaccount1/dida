package external

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"golang.org/x/oauth2"
	
	"dida/domain/entities"
	"dida/domain/services"
	"dida/domain/errors"
)

// OAuthAuthService OAuth认证服务实现
type OAuthAuthService struct {
	config       *oauth2.Config
	clientID     string
	clientSecret string
	authRepo     TokenRepository
}

// TokenRepository 令牌仓库接口
type TokenRepository interface {
	GetAccessToken(ctx context.Context) (string, error)
	GetRefreshToken(ctx context.Context) (string, error)
	SaveTokens(ctx context.Context, accessToken, refreshToken string) error
}

// NewOAuthAuthService 创建OAuth认证服务
func NewOAuthAuthService(
	clientID, clientSecret, authURL, tokenURL, redirectURL string,
	scopes []string,
	authRepo TokenRepository,
) services.AuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	
	return &OAuthAuthService{
		config:       config,
		clientID:     clientID,
		clientSecret: clientSecret,
		authRepo:     authRepo,
	}
}

// GetAuthorizationURL 获取授权URL
func (s *OAuthAuthService) GetAuthorizationURL(ctx context.Context) (string, error) {
	state := fmt.Sprintf("state_%d", time.Now().Unix())
	url := s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url, nil
}

// ExchangeCode 使用授权码交换访问令牌
func (s *OAuthAuthService) ExchangeCode(ctx context.Context, code string) (*entities.User, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.ErrAuthenticationFailed
	}
	
	// 获取用户信息
	userInfo, err := s.fetchUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}
	
	user := &entities.User{
		ID:           userInfo.ID,
		Username:     userInfo.Username,
		Email:        userInfo.Email,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	return user, nil
}

// RefreshToken 刷新访问令牌
func (s *OAuthAuthService) RefreshToken(ctx context.Context, refreshToken string) (*entities.User, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}
	
	tokenSource := s.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, errors.ErrRefreshTokenInvalid
	}
	
	// 获取用户信息
	userInfo, err := s.fetchUserInfo(ctx, newToken.AccessToken)
	if err != nil {
		return nil, err
	}
	
	user := &entities.User{
		ID:           userInfo.ID,
		Username:     userInfo.Username,
		Email:        userInfo.Email,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		TokenExpiry:  newToken.Expiry,
		UpdatedAt:    time.Now(),
	}
	
	return user, nil
}

// ValidateToken 验证令牌
func (s *OAuthAuthService) ValidateToken(ctx context.Context, token string) (*entities.User, error) {
	userInfo, err := s.fetchUserInfo(ctx, token)
	if err != nil {
		return nil, errors.ErrTokenExpired
	}
	
	user := &entities.User{
		ID:          userInfo.ID,
		Username:    userInfo.Username,
		Email:       userInfo.Email,
		AccessToken: token,
		UpdatedAt:   time.Now(),
	}
	
	return user, nil
}

// Login 登录（获取当前用户）
func (s *OAuthAuthService) Login(ctx context.Context) (*entities.User, error) {
	token, err := s.authRepo.GetAccessToken(ctx)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	
	return s.ValidateToken(ctx, token)
}

// Logout 登出
func (s *OAuthAuthService) Logout(ctx context.Context) error {
	// 清除本地令牌
	return s.authRepo.SaveTokens(ctx, "", "")
}

// GetCurrentUser 获取当前用户
func (s *OAuthAuthService) GetCurrentUser(ctx context.Context) (*entities.User, error) {
	return s.Login(ctx)
}

// IsAuthenticated 检查是否已认证
func (s *OAuthAuthService) IsAuthenticated(ctx context.Context) (bool, error) {
	token, err := s.authRepo.GetAccessToken(ctx)
	if err != nil {
		return false, nil
	}
	
	if token == "" {
		return false, nil
	}
	
	// 验证令牌是否有效
	_, err = s.ValidateToken(ctx, token)
	return err == nil, nil
}

// fetchUserInfo 获取用户信息
func (s *OAuthAuthService) fetchUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	// 这里需要根据TickTick的实际API来实现
	// 暂时使用模拟数据
	userInfo := &UserInfo{
		ID:       "user_id",
		Username: "user",
		Email:    "user@example.com",
	}
	
	return userInfo, nil
}

// exchangeCodeForToken 使用授权码交换令牌（手动实现，因为TickTick可能不完全兼容标准OAuth2）
func (s *OAuthAuthService) exchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.config.RedirectURL)
	
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	
	// 设置Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, errors.ErrAuthenticationFailed
	}
	
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	
	token := &oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
	}
	
	if tokenResp.ExpiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	
	return token, nil
}

// UserInfo 用户信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}