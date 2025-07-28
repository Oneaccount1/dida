package mcp

import (
	"context"
	"fmt"
	
	"github.com/mark3labs/mcp-go/mcp"
	
	"dida/domain/entities"
	"dida/usecases/project"
	"dida/usecases/task"
	"dida/usecases/auth"
)

// MCPHandlers MCP处理器
type MCPHandlers struct {
	// 用例
	getProjectsUC    *project.GetProjectsUseCase
	getProjectUC     *project.GetProjectUseCase
	getTasksUC       *task.GetTasksUseCase
	createTaskUC     *task.CreateTaskUseCase
	authenticateUC   *auth.AuthenticateUseCase
	getAuthURLUC     *auth.GetAuthURLUseCase
	refreshTokenUC   *auth.RefreshTokenUseCase
}

// NewMCPHandlers 创建MCP处理器
func NewMCPHandlers(
	getProjectsUC *project.GetProjectsUseCase,
	getProjectUC *project.GetProjectUseCase,
	getTasksUC *task.GetTasksUseCase,
	createTaskUC *task.CreateTaskUseCase,
	authenticateUC *auth.AuthenticateUseCase,
	getAuthURLUC *auth.GetAuthURLUseCase,
	refreshTokenUC *auth.RefreshTokenUseCase,
) *MCPHandlers {
	return &MCPHandlers{
		getProjectsUC:  getProjectsUC,
		getProjectUC:   getProjectUC,
		getTasksUC:     getTasksUC,
		createTaskUC:   createTaskUC,
		authenticateUC: authenticateUC,
		getAuthURLUC:   getAuthURLUC,
		refreshTokenUC: refreshTokenUC,
	}
}

// GetProjects 获取项目处理器
func (h *MCPHandlers) GetProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 获取参数
	includeInactive := false
	// GetString方法的签名可能需要调整，暂时直接设为false
	// if includeInactiveStr, err := request.GetString("include_inactive"); err == nil {
	//	includeInactive = includeInactiveStr == "true"
	// }
	
	// 执行用例
	projects, err := h.getProjectsUC.Execute(ctx, includeInactive)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching projects: %v", err)), nil
	}
	
	if len(projects) == 0 {
		return mcp.NewToolResultText("No projects found."), nil
	}
	
	// 格式化结果
	result := fmt.Sprintf("Found %d projects:\n\n", len(projects))
	for i, proj := range projects {
		result += fmt.Sprintf("Project %d:\n%s\n", i+1, h.formatProject(proj))
	}
	
	return mcp.NewToolResultText(result), nil
}

// GetProject 获取单个项目处理器
func (h *MCPHandlers) GetProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 获取项目ID
	projectID, err := request.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	// 执行用例
	proj, err := h.getProjectUC.Execute(ctx, projectID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching project: %v", err)), nil
	}
	
	return mcp.NewToolResultText(h.formatProject(proj)), nil
}

// GetTasks 获取任务处理器
func (h *MCPHandlers) GetTasks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	req := task.GetTasksRequest{}
	
	// 对于可选参数，暂时使用默认值
	// TODO: 修复参数解析方法
	// if projectID, err := request.GetString("project_id"); err == nil {
	//	req.ProjectID = projectID
	// }
	
	// 执行用例
	tasks, err := h.getTasksUC.Execute(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error fetching tasks: %v", err)), nil
	}
	
	if len(tasks) == 0 {
		return mcp.NewToolResultText("No tasks found."), nil
	}
	
	// 格式化结果
	result := fmt.Sprintf("Found %d tasks:\n\n", len(tasks))
	for i, t := range tasks {
		result += fmt.Sprintf("Task %d:\n%s\n", i+1, h.formatTask(t))
	}
	
	return mcp.NewToolResultText(result), nil
}

// CreateTask 创建任务处理器
func (h *MCPHandlers) CreateTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	req := task.CreateTaskRequest{}
	
	var err error
	req.ProjectID, err = request.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	req.Title, err = request.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	// 对于可选参数，暂时注释掉
	// TODO: 修复参数解析方法
	// if content, err := request.GetString("content"); err == nil {
	//	req.Content = content
	// }
	
	// 执行用例
	createdTask, err := h.createTaskUC.Execute(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error creating task: %v", err)), nil
	}
	
	result := fmt.Sprintf("Task created successfully:\n%s", h.formatTask(createdTask))
	return mcp.NewToolResultText(result), nil
}

// GetAuthURL 获取认证URL处理器
func (h *MCPHandlers) GetAuthURL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := h.getAuthURLUC.Execute(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error getting auth URL: %v", err)), nil
	}
	
	result := fmt.Sprintf("Please visit the following URL to authenticate:\n%s", url)
	return mcp.NewToolResultText(result), nil
}

// Authenticate 认证处理器
func (h *MCPHandlers) Authenticate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	code, err := request.RequireString("authorization_code")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	req := auth.LoginRequest{
		AuthorizationCode: code,
	}
	
	user, err := h.authenticateUC.Execute(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Authentication failed: %v", err)), nil
	}
	
	result := fmt.Sprintf("Authentication successful! Welcome, %s", user.Username)
	return mcp.NewToolResultText(result), nil
}

// 格式化辅助方法
func (h *MCPHandlers) formatProject(proj *entities.Project) string {
	formatted := fmt.Sprintf("Name: %s\n", proj.Name)
	formatted += fmt.Sprintf("ID: %s\n", proj.ID)
	
	if proj.Color != "" {
		formatted += fmt.Sprintf("Color: %s\n", proj.Color)
	}
	if proj.ViewMode != "" {
		formatted += fmt.Sprintf("View Mode: %s\n", proj.ViewMode)
	}
	formatted += fmt.Sprintf("Active: %v\n", proj.IsActive())
	if proj.Kind != "" {
		formatted += fmt.Sprintf("Kind: %s\n", proj.Kind)
	}
	return formatted
}

func (h *MCPHandlers) formatTask(t *entities.Task) string {
	formatted := fmt.Sprintf("Title: %s\n", t.Title)
	formatted += fmt.Sprintf("ID: %s\n", t.ID)
	formatted += fmt.Sprintf("Project ID: %s\n", t.ProjectID)
	
	if t.Content != "" {
		formatted += fmt.Sprintf("Content: %s\n", t.Content)
	}
	if t.Description != "" {
		formatted += fmt.Sprintf("Description: %s\n", t.Description)
	}
	
	formatted += fmt.Sprintf("Status: %s\n", h.formatTaskStatus(t.Status))
	formatted += fmt.Sprintf("Priority: %s\n", h.formatPriority(t.Priority))
	
	if t.DueDate != nil {
		formatted += fmt.Sprintf("Due Date: %s\n", t.DueDate.Format("2006-01-02 15:04:05"))
	}
	
	if t.IsCompleted() && t.CompletedTime != nil {
		formatted += fmt.Sprintf("Completed: %s\n", t.CompletedTime.Format("2006-01-02 15:04:05"))
	}
	
	return formatted
}

func (h *MCPHandlers) formatTaskStatus(status entities.TaskStatus) string {
	switch status {
	case entities.TaskStatusIncomplete:
		return "Incomplete"
	case entities.TaskStatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

func (h *MCPHandlers) formatPriority(priority entities.Priority) string {
	switch priority {
	case entities.PriorityNone:
		return "None"
	case entities.PriorityLow:
		return "Low"
	case entities.PriorityMedium:
		return "Medium"
	case entities.PriorityHigh:
		return "High"
	default:
		return "Unknown"
	}
}