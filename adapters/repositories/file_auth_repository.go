package repositories

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
	
	"dida/domain/entities"
	"dida/domain/repositories"
	"dida/domain/errors"
)

// FileAuthRepository 基于文件的认证仓库实现
type FileAuthRepository struct {
	filePath string
}

// AuthData 存储的认证数据结构
type AuthData struct {
	User         *entities.User `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// NewFileAuthRepository 创建文件认证仓库
func NewFileAuthRepository(filePath string) repositories.AuthRepository {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0755)
	
	return &FileAuthRepository{
		filePath: filePath,
	}
}

// GetUser 获取用户信息
func (r *FileAuthRepository) GetUser(ctx context.Context) (*entities.User, error) {
	data, err := r.loadAuthData()
	if err != nil {
		return nil, err
	}
	
	if data.User == nil {
		return nil, errors.ErrUnauthorized
	}
	
	return data.User, nil
}

// SaveUser 保存用户信息
func (r *FileAuthRepository) SaveUser(ctx context.Context, user *entities.User) error {
	data, _ := r.loadAuthData()
	data.User = user
	data.UpdatedAt = time.Now()
	
	return r.saveAuthData(data)
}

// DeleteUser 删除用户信息
func (r *FileAuthRepository) DeleteUser(ctx context.Context) error {
	return os.Remove(r.filePath)
}

// GetAccessToken 获取访问令牌
func (r *FileAuthRepository) GetAccessToken(ctx context.Context) (string, error) {
	data, err := r.loadAuthData()
	if err != nil {
		return "", err
	}
	
	if data.AccessToken == "" {
		return "", errors.ErrUnauthorized
	}
	
	return data.AccessToken, nil
}

// GetRefreshToken 获取刷新令牌
func (r *FileAuthRepository) GetRefreshToken(ctx context.Context) (string, error) {
	data, err := r.loadAuthData()
	if err != nil {
		return "", err
	}
	
	if data.RefreshToken == "" {
		return "", errors.ErrRefreshTokenInvalid
	}
	
	return data.RefreshToken, nil
}

// SaveTokens 保存令牌
func (r *FileAuthRepository) SaveTokens(ctx context.Context, accessToken, refreshToken string) error {
	data, _ := r.loadAuthData()
	data.AccessToken = accessToken
	data.RefreshToken = refreshToken
	data.UpdatedAt = time.Now()
	
	return r.saveAuthData(data)
}

// ClearTokens 清除令牌
func (r *FileAuthRepository) ClearTokens(ctx context.Context) error {
	data, _ := r.loadAuthData()
	data.AccessToken = ""
	data.RefreshToken = ""
	data.UpdatedAt = time.Now()
	
	return r.saveAuthData(data)
}

// IsAuthenticated 检查是否已认证
func (r *FileAuthRepository) IsAuthenticated(ctx context.Context) (bool, error) {
	data, err := r.loadAuthData()
	if err != nil {
		return false, nil // 文件不存在视为未认证
	}
	
	return data.User != nil && data.AccessToken != "" && !data.User.IsTokenExpired(), nil
}

// loadAuthData 加载认证数据
func (r *FileAuthRepository) loadAuthData() (*AuthData, error) {
	data := &AuthData{}
	
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return data, nil // 文件不存在返回空数据
	}
	
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(file, data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// saveAuthData 保存认证数据
func (r *FileAuthRepository) saveAuthData(data *AuthData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(r.filePath, jsonData, 0600)
}