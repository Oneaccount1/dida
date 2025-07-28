package auth

import (
	"context"
	"dida/domain/entities"
	"dida/domain/repositories"
	"dida/domain/services"
	"dida/domain/errors"
)

// AuthenticateUseCase 认证用例
type AuthenticateUseCase struct {
	authRepo repositories.AuthRepository
	authSvc  services.AuthService
}

// NewAuthenticateUseCase 创建认证用例
func NewAuthenticateUseCase(
	authRepo repositories.AuthRepository,
	authSvc services.AuthService,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		authRepo: authRepo,
		authSvc:  authSvc,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	AuthorizationCode string
}

// Execute 执行认证用例
func (uc *AuthenticateUseCase) Execute(ctx context.Context, req LoginRequest) (*entities.User, error) {
	// 1. 验证输入
	if req.AuthorizationCode == "" {
		return nil, errors.ErrRequiredField
	}

	// 2. 使用授权码交换访问令牌
	user, err := uc.authSvc.ExchangeCode(ctx, req.AuthorizationCode)
	if err != nil {
		return nil, errors.ErrAuthenticationFailed
	}

	// 3. 保存用户信息到仓库
	err = uc.authRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 4. 保存令牌到仓库
	err = uc.authRepo.SaveTokens(ctx, user.AccessToken, user.RefreshToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAuthURLUseCase 获取认证URL用例
type GetAuthURLUseCase struct {
	authSvc services.AuthService
}

// NewGetAuthURLUseCase 创建获取认证URL用例
func NewGetAuthURLUseCase(authSvc services.AuthService) *GetAuthURLUseCase {
	return &GetAuthURLUseCase{
		authSvc: authSvc,
	}
}

// Execute 执行获取认证URL用例
func (uc *GetAuthURLUseCase) Execute(ctx context.Context) (string, error) {
	return uc.authSvc.GetAuthorizationURL(ctx)
}

// RefreshTokenUseCase 刷新令牌用例
type RefreshTokenUseCase struct {
	authRepo repositories.AuthRepository
	authSvc  services.AuthService
}

// NewRefreshTokenUseCase 创建刷新令牌用例
func NewRefreshTokenUseCase(
	authRepo repositories.AuthRepository,
	authSvc services.AuthService,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		authRepo: authRepo,
		authSvc:  authSvc,
	}
}

// Execute 执行刷新令牌用例
func (uc *RefreshTokenUseCase) Execute(ctx context.Context) (*entities.User, error) {
	// 1. 获取当前刷新令牌
	refreshToken, err := uc.authRepo.GetRefreshToken(ctx)
	if err != nil {
		return nil, errors.ErrRefreshTokenInvalid
	}

	// 2. 使用刷新令牌获取新的访问令牌
	user, err := uc.authSvc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.ErrRefreshTokenInvalid
	}

	// 3. 更新用户信息
	err = uc.authRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 4. 更新令牌
	err = uc.authRepo.SaveTokens(ctx, user.AccessToken, user.RefreshToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}