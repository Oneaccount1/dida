package project

import (
	"context"
	"dida/domain/entities"
	"dida/domain/repositories"
	"dida/domain/services"
	"dida/domain/errors"
)

// GetProjectsUseCase 获取项目用例
type GetProjectsUseCase struct {
	projectRepo   repositories.ProjectRepository
	ticktickSvc   services.TickTickService
	authSvc       services.AuthService
}

// NewGetProjectsUseCase 创建获取项目用例
func NewGetProjectsUseCase(
	projectRepo repositories.ProjectRepository,
	ticktickSvc services.TickTickService,
	authSvc services.AuthService,
) *GetProjectsUseCase {
	return &GetProjectsUseCase{
		projectRepo: projectRepo,
		ticktickSvc: ticktickSvc,
		authSvc:     authSvc,
	}
}

// Execute 执行获取项目用例
func (uc *GetProjectsUseCase) Execute(ctx context.Context, includeInactive bool) ([]*entities.Project, error) {
	// 1. 检查用户认证状态
	authenticated, err := uc.authSvc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if !authenticated {
		return nil, errors.ErrUnauthorized
	}

	// 2. 从TickTick API获取项目
	projects, err := uc.ticktickSvc.GetProjects(ctx)
	if err != nil {
		return nil, err
	}

	// 3. 过滤不活跃的项目（如果需要）
	if !includeInactive {
		activeProjects := make([]*entities.Project, 0)
		for _, project := range projects {
			if project.IsActive() {
				activeProjects = append(activeProjects, project)
			}
		}
		projects = activeProjects
	}

	// 4. 可选：缓存到本地仓库
	// TODO: 实现缓存逻辑

	return projects, nil
}