package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	
	"dida/domain/entities"
	"dida/domain/services"
	"dida/domain/errors"
)

// TickTickClient TickTick API客户端实现
type TickTickClient struct {
	baseURL    string
	httpClient *http.Client
	authRepo   AuthTokenProvider
}

// AuthTokenProvider 提供认证令牌的接口
type AuthTokenProvider interface {
	GetAccessToken(ctx context.Context) (string, error)
}

// NewTickTickClient 创建TickTick客户端
func NewTickTickClient(baseURL string, authRepo AuthTokenProvider) services.TickTickService {
	return &TickTickClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		authRepo:   authRepo,
	}
}

// GetProjects 获取所有项目
func (c *TickTickClient) GetProjects(ctx context.Context) ([]*entities.Project, error) {
	var apiProjects []APIProject
	err := c.makeRequest(ctx, "GET", "/projects", nil, &apiProjects)
	if err != nil {
		return nil, err
	}
	
	projects := make([]*entities.Project, len(apiProjects))
	for i, apiProject := range apiProjects {
		projects[i] = c.convertAPIProjectToEntity(apiProject)
	}
	
	return projects, nil
}

// GetProject 获取单个项目
func (c *TickTickClient) GetProject(ctx context.Context, projectID string) (*entities.Project, error) {
	var apiProject APIProject
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/projects/%s", projectID), nil, &apiProject)
	if err != nil {
		return nil, err
	}
	
	return c.convertAPIProjectToEntity(apiProject), nil
}

// CreateProject 创建项目
func (c *TickTickClient) CreateProject(ctx context.Context, project *entities.Project) error {
	apiProject := c.convertEntityToAPIProject(project)
	var response APIProject
	err := c.makeRequest(ctx, "POST", "/projects", apiProject, &response)
	if err != nil {
		return err
	}
	
	// 更新实体ID
	project.ID = response.ID
	return nil
}

// UpdateProject 更新项目
func (c *TickTickClient) UpdateProject(ctx context.Context, project *entities.Project) error {
	apiProject := c.convertEntityToAPIProject(project)
	return c.makeRequest(ctx, "PUT", fmt.Sprintf("/projects/%s", project.ID), apiProject, nil)
}

// DeleteProject 删除项目
func (c *TickTickClient) DeleteProject(ctx context.Context, projectID string) error {
	return c.makeRequest(ctx, "DELETE", fmt.Sprintf("/projects/%s", projectID), nil, nil)
}

// GetTasks 获取任务
func (c *TickTickClient) GetTasks(ctx context.Context, projectID string) ([]*entities.Task, error) {
	var apiTasks []APITask
	url := fmt.Sprintf("/projects/%s/tasks", projectID)
	err := c.makeRequest(ctx, "GET", url, nil, &apiTasks)
	if err != nil {
		return nil, err
	}
	
	tasks := make([]*entities.Task, len(apiTasks))
	for i, apiTask := range apiTasks {
		tasks[i] = c.convertAPITaskToEntity(apiTask)
	}
	
	return tasks, nil
}

// GetTask 获取单个任务
func (c *TickTickClient) GetTask(ctx context.Context, taskID string) (*entities.Task, error) {
	var apiTask APITask
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/tasks/%s", taskID), nil, &apiTask)
	if err != nil {
		return nil, err
	}
	
	return c.convertAPITaskToEntity(apiTask), nil
}

// CreateTask 创建任务
func (c *TickTickClient) CreateTask(ctx context.Context, task *entities.Task) error {
	apiTask := c.convertEntityToAPITask(task)
	var response APITask
	err := c.makeRequest(ctx, "POST", "/tasks", apiTask, &response)
	if err != nil {
		return err
	}
	
	// 更新实体ID
	task.ID = response.ID
	return nil
}

// UpdateTask 更新任务
func (c *TickTickClient) UpdateTask(ctx context.Context, task *entities.Task) error {
	apiTask := c.convertEntityToAPITask(task)
	return c.makeRequest(ctx, "PUT", fmt.Sprintf("/tasks/%s", task.ID), apiTask, nil)
}

// DeleteTask 删除任务
func (c *TickTickClient) DeleteTask(ctx context.Context, taskID string) error {
	return c.makeRequest(ctx, "DELETE", fmt.Sprintf("/tasks/%s", taskID), nil, nil)
}

// CompleteTask 完成任务
func (c *TickTickClient) CompleteTask(ctx context.Context, taskID string) error {
	payload := map[string]interface{}{
		"status": 1, // 1 表示已完成
	}
	return c.makeRequest(ctx, "PATCH", fmt.Sprintf("/tasks/%s", taskID), payload, nil)
}

// SyncProjects 同步项目
func (c *TickTickClient) SyncProjects(ctx context.Context) ([]*entities.Project, error) {
	return c.GetProjects(ctx)
}

// SyncTasks 同步任务
func (c *TickTickClient) SyncTasks(ctx context.Context, projectID string) ([]*entities.Task, error) {
	return c.GetTasks(ctx, projectID)
}

// HealthCheck 健康检查
func (c *TickTickClient) HealthCheck(ctx context.Context) error {
	return c.makeRequest(ctx, "GET", "/health", nil, nil)
}

// makeRequest 通用HTTP请求方法
func (c *TickTickClient) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(jsonData)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}
	
	// 添加认证头
	token, err := c.authRepo.GetAccessToken(ctx)
	if err != nil {
		return errors.ErrUnauthorized
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.ErrNetworkFailure
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 401 {
		return errors.ErrUnauthorized
	}
	if resp.StatusCode == 404 {
		if strings.Contains(path, "/projects/") {
			return errors.ErrProjectNotFound
		}
		if strings.Contains(path, "/tasks/") {
			return errors.ErrTaskNotFound
		}
	}
	if resp.StatusCode >= 400 {
		return errors.ErrServiceUnavailable
	}
	
	if result != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		
		return json.Unmarshal(respBody, result)
	}
	
	return nil
}

// API模型转换方法
func (c *TickTickClient) convertAPIProjectToEntity(apiProject APIProject) *entities.Project {
	return &entities.Project{
		ID:       apiProject.ID,
		Name:     apiProject.Name,
		Color:    apiProject.Color,
		ViewMode: entities.ViewMode(apiProject.ViewMode),
		Closed:   apiProject.Closed,
		Kind:     entities.ProjectKind(apiProject.Kind),
	}
}

func (c *TickTickClient) convertEntityToAPIProject(project *entities.Project) APIProject {
	return APIProject{
		ID:       project.ID,
		Name:     project.Name,
		Color:    project.Color,
		ViewMode: string(project.ViewMode),
		Closed:   project.Closed,
		Kind:     string(project.Kind),
	}
}

func (c *TickTickClient) convertAPITaskToEntity(apiTask APITask) *entities.Task {
	var startDate, dueDate, completedTime *time.Time
	
	if apiTask.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, apiTask.StartDate); err == nil {
			startDate = &t
		}
	}
	if apiTask.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, apiTask.DueDate); err == nil {
			dueDate = &t
		}
	}
	if apiTask.CompletedTime != "" {
		if t, err := time.Parse(time.RFC3339, apiTask.CompletedTime); err == nil {
			completedTime = &t
		}
	}
	
	return &entities.Task{
		ID:            apiTask.ID,
		ProjectID:     apiTask.ProjectID,
		Title:         apiTask.Title,
		Content:       apiTask.Content,
		Description:   apiTask.Desc,
		IsAllDay:      apiTask.IsAllDay,
		StartDate:     startDate,
		DueDate:       dueDate,
		TimeZone:      apiTask.TimeZone,
		Reminders:     apiTask.Reminders,
		RepeatFlag:    apiTask.RepeatFlag,
		Priority:      entities.Priority(apiTask.Priority),
		Status:        entities.TaskStatus(apiTask.Status),
		CompletedTime: completedTime,
		SortOrder:     apiTask.SortOrder,
	}
}

func (c *TickTickClient) convertEntityToAPITask(task *entities.Task) APITask {
	apiTask := APITask{
		ID:         task.ID,
		ProjectID:  task.ProjectID,
		Title:      task.Title,
		Content:    task.Content,
		Desc:       task.Description,
		IsAllDay:   task.IsAllDay,
		TimeZone:   task.TimeZone,
		Reminders:  task.Reminders,
		RepeatFlag: task.RepeatFlag,
		Priority:   int(task.Priority),
		Status:     int(task.Status),
		SortOrder:  task.SortOrder,
	}
	
	if task.StartDate != nil {
		apiTask.StartDate = task.StartDate.Format(time.RFC3339)
	}
	if task.DueDate != nil {
		apiTask.DueDate = task.DueDate.Format(time.RFC3339)
	}
	if task.CompletedTime != nil {
		apiTask.CompletedTime = task.CompletedTime.Format(time.RFC3339)
	}
	
	return apiTask
}

// API模型定义
type APIProject struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Color    string `json:"color,omitempty"`
	ViewMode string `json:"viewMode,omitempty"`
	Closed   bool   `json:"closed,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

type APITask struct {
	ID            string   `json:"id,omitempty"`
	ProjectID     string   `json:"projectId"`
	Title         string   `json:"title"`
	Content       string   `json:"content,omitempty"`
	Desc          string   `json:"desc,omitempty"`
	IsAllDay      bool     `json:"isAllDay,omitempty"`
	StartDate     string   `json:"startDate,omitempty"`
	DueDate       string   `json:"dueDate,omitempty"`
	TimeZone      string   `json:"timeZone,omitempty"`
	Reminders     []string `json:"reminders,omitempty"`
	RepeatFlag    string   `json:"repeatFlag,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	Status        int      `json:"status,omitempty"`
	CompletedTime string   `json:"completedTime,omitempty"`
	SortOrder     int      `json:"sortOrder,omitempty"`
}