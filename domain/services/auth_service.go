package services

import (
	"context"
	"dida/domain/entities"
)

// AuthService 定义认证服务接口
type AuthService interface {
	// OAuth流程
	GetAuthorizationURL(ctx context.Context) (string, error)
	ExchangeCode(ctx context.Context, code string) (*entities.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entities.User, error)
	
	// 令牌验证
	ValidateToken(ctx context.Context, token string) (*entities.User, error)
	
	// 登录登出
	Login(ctx context.Context) (*entities.User, error)
	Logout(ctx context.Context) error
	
	// 令牌管理
	GetCurrentUser(ctx context.Context) (*entities.User, error)
	IsAuthenticated(ctx context.Context) (bool, error)
}