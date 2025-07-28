package repositories

import (
	"context"
	"dida/domain/entities"
)

// ProjectRepository 定义项目数据访问接口
type ProjectRepository interface {
	// 获取项目
	GetByID(ctx context.Context, projectID string) (*entities.Project, error)
	GetAll(ctx context.Context) ([]*entities.Project, error)
	GetActiveProjects(ctx context.Context) ([]*entities.Project, error)
	GetProjectsByKind(ctx context.Context, kind entities.ProjectKind) ([]*entities.Project, error)
	
	// 创建和更新项目
	Create(ctx context.Context, project *entities.Project) error
	Update(ctx context.Context, project *entities.Project) error
	Delete(ctx context.Context, projectID string) error
	
	// 项目状态操作
	Close(ctx context.Context, projectID string) error
	Reopen(ctx context.Context, projectID string) error
}