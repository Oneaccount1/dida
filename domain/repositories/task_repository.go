package repositories

import (
	"context"
	"dida/domain/entities"
)

// TaskRepository 定义任务数据访问接口
type TaskRepository interface {
	// 获取任务
	GetByID(ctx context.Context, taskID string) (*entities.Task, error)
	GetByProjectID(ctx context.Context, projectID string) ([]*entities.Task, error)
	GetAll(ctx context.Context) ([]*entities.Task, error)
	
	// 创建和更新任务
	Create(ctx context.Context, task *entities.Task) error
	Update(ctx context.Context, task *entities.Task) error
	Delete(ctx context.Context, taskID string) error
	
	// 任务状态操作
	MarkComplete(ctx context.Context, taskID string) error
	MarkIncomplete(ctx context.Context, taskID string) error
	
	// 查询操作
	GetOverdueTasks(ctx context.Context) ([]*entities.Task, error)
	GetTasksByPriority(ctx context.Context, priority entities.Priority) ([]*entities.Task, error)
	GetTasksByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error)
	
	// 批量操作
	CreateBatch(ctx context.Context, tasks []*entities.Task) error
	UpdateBatch(ctx context.Context, tasks []*entities.Task) error
}