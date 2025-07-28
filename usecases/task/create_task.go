package task

import (
	"context"
	"time"
	"dida/domain/entities"
	"dida/domain/repositories"
	"dida/domain/services"
	"dida/domain/errors"
)

// CreateTaskUseCase 创建任务用例
type CreateTaskUseCase struct {
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	ticktickSvc services.TickTickService
	authSvc     services.AuthService
}

// NewCreateTaskUseCase 创建任务用例构造函数
func NewCreateTaskUseCase(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	ticktickSvc services.TickTickService,
	authSvc services.AuthService,
) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		ticktickSvc: ticktickSvc,
		authSvc:     authSvc,
	}
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	ProjectID   string
	Title       string
	Content     string
	Description string
	IsAllDay    bool
	StartDate   *time.Time
	DueDate     *time.Time
	TimeZone    string
	Reminders   []string
	RepeatFlag  string
	Priority    entities.Priority
}

// Execute 执行创建任务用例
func (uc *CreateTaskUseCase) Execute(ctx context.Context, req CreateTaskRequest) (*entities.Task, error) {
	// 1. 检查用户认证状态
	authenticated, err := uc.authSvc.IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if !authenticated {
		return nil, errors.ErrUnauthorized
	}

	// 2. 验证必填字段
	if req.Title == "" {
		return nil, errors.ErrRequiredField
	}
	if req.ProjectID == "" {
		return nil, errors.ErrRequiredField
	}

	// 3. 验证项目存在且活跃
	project, err := uc.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, errors.ErrProjectNotFound
	}
	if !project.IsActive() {
		return nil, errors.ErrProjectClosed
	}

	// 4. 创建任务实体
	now := time.Now()
	task := &entities.Task{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		IsAllDay:    req.IsAllDay,
		StartDate:   req.StartDate,
		DueDate:     req.DueDate,
		TimeZone:    req.TimeZone,
		Reminders:   req.Reminders,
		RepeatFlag:  req.RepeatFlag,
		Priority:    req.Priority,
		Status:      entities.TaskStatusIncomplete,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 5. 通过TickTick API创建任务
	err = uc.ticktickSvc.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	// 6. 可选：缓存到本地仓库
	_ = uc.taskRepo.Create(ctx, task)

	return task, nil
}