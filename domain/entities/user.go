package entities

import "time"

// User 表示TickTick用户实体
type User struct {
	ID           string
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Business methods for User entity

// IsTokenExpired 检查访问令牌是否已过期
func (u *User) IsTokenExpired() bool {
	return time.Now().After(u.TokenExpiry)
}

// UpdateTokens 更新用户令牌
func (u *User) UpdateTokens(accessToken, refreshToken string, expiry time.Time) {
	u.AccessToken = accessToken
	u.RefreshToken = refreshToken
	u.TokenExpiry = expiry
	u.UpdatedAt = time.Now()
}

// IsAuthenticated 检查用户是否已认证
func (u *User) IsAuthenticated() bool {
	return u.AccessToken != "" && !u.IsTokenExpired()
}