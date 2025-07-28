package services

import (
	"context"
	"dida/domain/entities"
)

// TickTickService 定义与TickTick API交互的服务接口
type TickTickService interface {
	// 项目操作
	GetProjects(ctx context.Context) ([]*entities.Project, error)
	GetProject(ctx context.Context, projectID string) (*entities.Project, error)
	CreateProject(ctx context.Context, project *entities.Project) error
	UpdateProject(ctx context.Context, project *entities.Project) error
	DeleteProject(ctx context.Context, projectID string) error
	
	// 任务操作
	GetTasks(ctx context.Context, projectID string) ([]*entities.Task, error)
	GetTask(ctx context.Context, taskID string) (*entities.Task, error)
	CreateTask(ctx context.Context, task *entities.Task) error
	UpdateTask(ctx context.Context, task *entities.Task) error
	DeleteTask(ctx context.Context, taskID string) error
	CompleteTask(ctx context.Context, taskID string) error
	
	// 批量操作
	SyncProjects(ctx context.Context) ([]*entities.Project, error)
	SyncTasks(ctx context.Context, projectID string) ([]*entities.Task, error)
	
	// 健康检查
	HealthCheck(ctx context.Context) error
}