package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dida/internal/auth"
	"dida/internal/config"
	"dida/internal/errors"
)

// TickTickClient 定义了与TickTick API交互的客户端
type TickTickClient struct {
	config     *config.Config
	HTTPClient *http.Client
	auth       *auth.TickTickAuth
}

func NewTickTickClient() (*TickTickClient, error) {

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout: cfg.TickTick.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	// 创建认证管理器
	tickAuth, err := auth.NewTickTickAuth(cfg.TickTick.ClientID, cfg.TickTick.ClientSecret)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrClientInit, err, "failed to initialize auth manager")
	}

	return &TickTickClient{
		config:     cfg,
		HTTPClient: httpClient,
		auth:       tickAuth,
	}, nil
}
func (c *TickTickClient) RefreshAccessToken() error {
	if c.config.TickTick.RefreshToken == "" {
		return errors.New(errors.ErrTokenRefreshFailed, "no refresh token available")
	}

	// 使用 auth 包进行令牌刷新
	tokens, err := c.auth.RefreshAccessToken(c.config.TickTick.RefreshToken)
	if err != nil {
		return err
	}

	// 更新配置中的令牌
	c.config.TickTick.AccessToken = tokens.AccessToken
	if tokens.RefreshToken != "" {
		c.config.TickTick.RefreshToken = tokens.RefreshToken
	}

	return nil
}

// GetAccessToken 获取当前的访问令牌
func (c *TickTickClient) GetAccessToken() string {
	return c.config.TickTick.AccessToken
}
func (c *TickTickClient) makeRequest(method, endpoint string, data interface{}) ([]byte, error) {
	// 直接执行请求
	return c.doRequest(method, endpoint, data)
}

func (c *TickTickClient) doRequest(method, endpoint string, data interface{}) ([]byte, error) {
	// 构建完整URL
	url := c.config.TickTick.BaseURL + endpoint

	// 构建请求体
	var reqBody io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrapf(errors.ErrDataMarshal, err, "failed to marshal request data")
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// 创建请求
	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrAPIRequest, err, "failed to create HTTP request")
	}

	// 设置请求头
	request.Header.Set("Authorization", "Bearer "+c.config.TickTick.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept-Encoding", "identity")
	request.Header.Set("User-Agent", "ticktick-mcp-go/1.0")

	// 发送请求
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrAPIRequest, err, "failed to send HTTP request")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrAPIResponse, err, "failed to read response body")
	}

	// 检查是否需要刷新令牌（如果是401未授权）
	if response.StatusCode == http.StatusUnauthorized {
		// 检查是否有 refresh token
		if c.config.TickTick.RefreshToken == "" {
			return nil, errors.New(errors.ErrTokenRefreshFailed,
				"Access token expired and no refresh token available. Please use the oauth_authorize tool to re-authenticate.")
		}

		if err := c.RefreshAccessToken(); err != nil {
			return nil, errors.Wrapf(errors.ErrTokenRefreshFailed, err, "failed to refresh access token")
		}

		// 重试请求
		request.Header.Set("Authorization", "Bearer "+c.config.TickTick.AccessToken)
		response, err = c.HTTPClient.Do(request)
		if err != nil {
			return nil, errors.Wrapf(errors.ErrAPIRequest, err, "failed to send request after token refresh")
		}
		defer response.Body.Close()

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrapf(errors.ErrAPIResponse, err, "failed to read response after token refresh")
		}
	}

	// 检查HTTP状态码
	if response.StatusCode >= 400 {
		return nil, errors.Newf(errors.ErrAPIResponse, "API error %s: %s", response.Status, string(body))
	}

	return body, nil
}

// GetProjects 获取所有项目
func (c *TickTickClient) GetProjects() ([]Project, error) {
	body, err := c.makeRequest("GET", "/project", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err = json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("error unmarshalling project: %v", err)
	}

	return projects, nil
}

// GetProject 获取特定项目
func (c *TickTickClient) GetProject(projectID string) (*Project, error) {
	body, err := c.makeRequest("GET", "/project/"+projectID, nil)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshalling project: %v", err)
	}
	return &project, nil
}
func (c *TickTickClient) GetProjectWithData(projectID string) (*ProjectData, error) {
	body, err := c.makeRequest("GET", "/project/"+projectID+"/data", nil)
	if err != nil {
		return nil, err
	}
	var projectData ProjectData
	if err := json.Unmarshal(body, &projectData); err != nil {
		return nil, fmt.Errorf("error unmarshlling project data: %v", err)
	}
	return &projectData, nil
}
func (c *TickTickClient) CreateProject(project Project) (*Project, error) {
	body, err := c.makeRequest("POST", "/project", project)
	if err != nil {
		return nil, err
	}

	var createdProject Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, err
	}
	return &createdProject, nil
}
func (c *TickTickClient) UpdateProject(project Project) (*Project, error) {
	body, err := c.makeRequest("POST", "/project/"+project.ID, project)
	if err != nil {
		return nil, err
	}

	var updatedProject Project
	if err := json.Unmarshal(body, &updatedProject); err != nil {
		return nil, err
	}
	return &updatedProject, nil

}
func (c *TickTickClient) DeleteProject(projectID string) error {
	_, err := c.makeRequest("DELETE", "/project/"+projectID, "")
	if err != nil {
		return err
	}
	return nil
}

// CreateTask 创建任务
func (c *TickTickClient) CreateTask(task Task) (*Task, error) {
	body, err := c.makeRequest("POST", "/task", task)
	if err != nil {
		return nil, err
	}
	var createdTask Task
	if err := json.Unmarshal(body, &createdTask); err != nil {
		return nil, fmt.Errorf("error unmarshalling created task: %v", err)
	}
	return &createdTask, nil
}

// UpdateTask 更新任务
func (c *TickTickClient) UpdateTask(task Task) (*Task, error) {
	body, err := c.makeRequest("POST", "/task/"+task.ID, task)
	if err != nil {
		return nil, err
	}
	var updatedTask Task

	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("error unmarshalling updated task: %v", err)
	}
	return &updatedTask, nil
}

// CompletedTask 完成任务
func (c *TickTickClient) CompletedTask(projectID, taskID string) error {
	_, err := c.makeRequest("POST", "/project/"+projectID+"/task/"+taskID+"/complete", "")
	if err != nil {
		return err
	}
	return nil
}

// DeleteTask 删除任务
func (c *TickTickClient) DeleteTask(projectID, taskID string) error {
	_, err := c.makeRequest("DELETE", "/project/"+projectID+"/task/"+taskID, "")
	if err != nil {
		return err
	}
	return nil
}

// GetTask 获取特定任务
func (c *TickTickClient) GetTask(projectID, taskID string) (*Task, error) {
	body, err := c.makeRequest("GET", fmt.Sprintf("/project/%s/task/%s", projectID, taskID), nil)
	if err != nil {
		return nil, err
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("error unmarshalling task: %v", err)
	}

	return &task, nil
}
