package repositories

import (
	"context"
	"dida/domain/entities"
)

// AuthRepository 定义认证数据访问接口
type AuthRepository interface {
	// 用户管理
	GetUser(ctx context.Context) (*entities.User, error)
	SaveUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context) error
	
	// 令牌管理
	GetAccessToken(ctx context.Context) (string, error)
	GetRefreshToken(ctx context.Context) (string, error)
	SaveTokens(ctx context.Context, accessToken, refreshToken string) error
	ClearTokens(ctx context.Context) error
	
	// 认证状态
	IsAuthenticated(ctx context.Context) (bool, error)
}