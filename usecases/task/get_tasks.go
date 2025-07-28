package task

import (
	"context"
	"dida/domain/entities"
	"dida/domain/repositories"
	"dida/domain/services"
	"dida/domain/errors"
)

// GetTasksUseCase 获取任务用例
type GetTasksUseCase struct {
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	ticktickSvc services.TickTickService
	authSvc     services.AuthService
}

// NewGetTasksUseCase 创建获取任务用例
func NewGetTasksUseCase(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	ticktickSvc services.TickTickService,
	authSvc services.AuthService,
) *GetTasksUseCase {
	return &GetTasksUseCase{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		ticktickSvc: ticktickSvc,
		authSvc:     authSvc,
	}
}

// GetTasksRequest 获取任务请求
type GetTasksRequest struct {
	ProjectID      string
	IncludeCompleted bool
	Priority       *entities.Priority
	Status         *entities.TaskStatus
}

// Execute 执行获取任务用例
func (uc *GetTasksUseCase) Execute(ctx context.Context, req GetTasksRequest) ([]*entities.Task, error) {
	// 1. 检查用户认证状态
	authenticated, err := uc.authSvc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if !authenticated {
		return nil, errors.ErrUnauthorized
	}

	// 2. 验证项目存在
	if req.ProjectID != "" {
		_, err := uc.projectRepo.GetByID(ctx, req.ProjectID)
		if err != nil {
			return nil, errors.ErrProjectNotFound
		}
	}

	// 3. 获取任务
	var tasks []*entities.Task
	if req.ProjectID != "" {
		tasks, err = uc.ticktickSvc.GetTasks(ctx, req.ProjectID)
	} else {
		tasks, err = uc.taskRepo.GetAll(ctx)
	}
	if err != nil {
		return nil, err
	}

	// 4. 应用过滤条件
	filteredTasks := make([]*entities.Task, 0)
	for _, task := range tasks {
		// 过滤已完成任务
		if !req.IncludeCompleted && task.IsCompleted() {
			continue
		}

		// 过滤优先级
		if req.Priority != nil && task.Priority != *req.Priority {
			continue
		}

		// 过滤状态
		if req.Status != nil && task.Status != *req.Status {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, nil
}