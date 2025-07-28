package app

import (
	"dida/adapters/external"
	"dida/adapters/repositories"
	"dida/infrastructure/config"
	"dida/interfaces/mcp"
	"dida/usecases/auth"
	"dida/usecases/project"
	"dida/usecases/task"
)

// Application 应用程序
type Application struct {
	Config    *config.Config
	MCPServer *mcp.MCPServer
}

// NewApplication 创建应用程序
func NewApplication() (*Application, error) {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	
	// 2. 创建仓库层
	authRepo := repositories.NewFileAuthRepository(cfg.Auth.TokenFilePath)
	
	// 3. 创建外部服务层
	ticktickService := external.NewTickTickClient(cfg.TickTick.BaseURL, authRepo)
	authService := external.NewOAuthAuthService(
		cfg.TickTick.ClientID,
		cfg.TickTick.ClientSecret,
		cfg.TickTick.AuthURL,
		cfg.TickTick.TokenURL,
		cfg.TickTick.RedirectURL,
		cfg.TickTick.Scopes,
		authRepo,
	)
	
	// 4. 创建用例层
	// 项目用例
	getProjectsUC := project.NewGetProjectsUseCase(nil, ticktickService, authService)
	getProjectUC := project.NewGetProjectUseCase(nil, ticktickService, authService)
	
	// 任务用例
	getTasksUC := task.NewGetTasksUseCase(nil, nil, ticktickService, authService)
	createTaskUC := task.NewCreateTaskUseCase(nil, nil, ticktickService, authService)
	
	// 认证用例
	authenticateUC := auth.NewAuthenticateUseCase(authRepo, authService)
	getAuthURLUC := auth.NewGetAuthURLUseCase(authService)
	refreshTokenUC := auth.NewRefreshTokenUseCase(authRepo, authService)
	
	// 5. 创建接口层
	handlers := mcp.NewMCPHandlers(
		getProjectsUC,
		getProjectUC,
		getTasksUC,
		createTaskUC,
		authenticateUC,
		getAuthURLUC,
		refreshTokenUC,
	)
	
	mcpServer := mcp.NewMCPServer(cfg.Server.Name, cfg.Server.Version, handlers)
	
	return &Application{
		Config:    cfg,
		MCPServer: mcpServer,
	}, nil
}

// Start 启动应用程序
func (app *Application) Start() error {
	return app.MCPServer.Start()
}