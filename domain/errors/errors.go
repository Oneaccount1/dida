package errors

import "errors"

// 认证相关错误
var (
	ErrUnauthorized         = errors.New("user not authenticated")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenExpired         = errors.New("access token expired")
	ErrRefreshTokenInvalid  = errors.New("refresh token invalid")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

// 项目相关错误
var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectExists       = errors.New("project already exists")
	ErrProjectClosed       = errors.New("project is closed")
	ErrInvalidProjectData  = errors.New("invalid project data")
)

// 任务相关错误
var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskExists         = errors.New("task already exists")
	ErrTaskAlreadyComplete = errors.New("task already completed")
	ErrInvalidTaskData    = errors.New("invalid task data")
	ErrTaskDependency     = errors.New("task has dependencies")
)

// 网络和外部服务错误
var (
	ErrNetworkFailure     = errors.New("network failure")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrAPIQuotaExceeded   = errors.New("API quota exceeded")
	ErrInvalidResponse    = errors.New("invalid response from service")
)

// 配置相关错误
var (
	ErrConfigNotFound     = errors.New("configuration not found")
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrMissingCredentials = errors.New("missing API credentials")
)

// 验证错误
var (
	ErrValidationFailed = errors.New("validation failed")
	ErrRequiredField    = errors.New("required field missing")
	ErrInvalidFormat    = errors.New("invalid format")
)