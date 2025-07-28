package project

import (
	"context"
	"dida/domain/entities"
	"dida/domain/repositories"
	"dida/domain/services"
	"dida/domain/errors"
)

// GetProjectUseCase 获取单个项目用例
type GetProjectUseCase struct {
	projectRepo repositories.ProjectRepository
	ticktickSvc services.TickTickService
	authSvc     services.AuthService
}

// NewGetProjectUseCase 创建获取项目用例
func NewGetProjectUseCase(
	projectRepo repositories.ProjectRepository,
	ticktickSvc services.TickTickService,
	authSvc services.AuthService,
) *GetProjectUseCase {
	return &GetProjectUseCase{
		projectRepo: projectRepo,
		ticktickSvc: ticktickSvc,
		authSvc:     authSvc,
	}
}

// Execute 执行获取项目用例
func (uc *GetProjectUseCase) Execute(ctx context.Context, projectID string) (*entities.Project, error) {
	// 1. 验证输入
	if projectID == "" {
		return nil, errors.ErrRequiredField
	}

	// 2. 检查用户认证状态
	authenticated, err := uc.authSvc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if !authenticated {
		return nil, errors.ErrUnauthorized
	}

	// 3. 先尝试从本地仓库获取
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err == nil {
		return project, nil
	}

	// 4. 本地没有，从TickTick API获取
	project, err = uc.ticktickSvc.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 5. 可选：缓存到本地仓库
	_ = uc.projectRepo.Create(ctx, project)

	return project, nil
}