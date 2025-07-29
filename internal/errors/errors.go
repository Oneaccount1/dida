package errors

import (
	"fmt"
)

// ErrorCode 定义错误代码类型
type ErrorCode string

const (
	// 认证相关错误
	ErrAuthFailed         ErrorCode = "AUTH_FAILED"
	ErrTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrTokenRefreshFailed ErrorCode = "TOKEN_REFRESH_FAILED"
	ErrInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"

	// 客户端相关错误
	ErrClientInit     ErrorCode = "CLIENT_INIT_FAILED"
	ErrAPIRequest     ErrorCode = "API_REQUEST_FAILED"
	ErrAPIResponse    ErrorCode = "API_RESPONSE_ERROR"
	ErrNetworkTimeout ErrorCode = "NETWORK_TIMEOUT"

	// 数据相关错误
	ErrInvalidData   ErrorCode = "INVALID_DATA"
	ErrDataNotFound  ErrorCode = "DATA_NOT_FOUND"
	ErrDataMarshal   ErrorCode = "DATA_MARSHAL_ERROR"
	ErrDataUnmarshal ErrorCode = "DATA_UNMARSHAL_ERROR"

	// 配置相关错误
	ErrConfigLoad   ErrorCode = "CONFIG_LOAD_FAILED"
	ErrEnvLoad      ErrorCode = "ENV_LOAD_FAILED"
	ErrFileNotFound ErrorCode = "FILE_NOT_FOUND"

	// 服务器相关错误
	ErrServerStart    ErrorCode = "SERVER_START_FAILED"
	ErrServerShutdown ErrorCode = "SERVER_SHUTDOWN_FAILED"
)

// AppError 应用程序错误结构
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New 创建新的应用程序错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装现有错误
func Wrap(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Newf 创建格式化的应用程序错误
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrapf 包装现有错误并格式化消息
func Wrapf(code ErrorCode, cause error, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// IsCode 检查错误是否为指定的错误代码
func IsCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// GetCode 获取错误代码
func GetCode(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ""
}
