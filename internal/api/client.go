package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TickTickClient 定义了与TickTick API交互的客户端
type TickTickClient struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	BaseURL      string
	TokenURL     string
	HTTPClient   *http.Client
}

// Task 表示TickTick任务
type Task struct {
	ID            string     `json:"id,omitempty"`
	ProjectID     string     `json:"projectId"`
	Title         string     `json:"title"`
	Content       string     `json:"content,omitempty"`
	Desc          string     `json:"desc,omitempty"`
	IsAllDay      bool       `json:"isAllDay,omitempty"`
	StartDate     string     `json:"startDate,omitempty"`
	DueDate       string     `json:"dueDate,omitempty"`
	TimeZone      string     `json:"timeZone,omitempty"`
	Reminders     []string   `json:"reminders,omitempty"`
	RepeatFlag    string     `json:"repeatFlag,omitempty"`
	Priority      int        `json:"priority,omitempty"`
	Status        int        `json:"status,omitempty"`
	CompletedTime string     `json:"completedTime,omitempty"`
	SortOrder     int        `json:"sortOrder,omitempty"`
	Items         []TaskItem `json:"items,omitempty"`
}

// TaskItem 表示子任务
type TaskItem struct {
	ID            string `json:"id,omitempty"`
	Status        int    `json:"status"`
	Title         string `json:"title"`
	SortOrder     int    `json:"sortOrder,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	IsAllDay      bool   `json:"isAllDay,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
	CompletedTime string `json:"completedTime,omitempty"`
}

// Project 表示TickTick项目
type Project struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Color    string `json:"color,omitempty"`
	ViewMode string `json:"viewMode,omitempty"`
	Closed   bool   `json:"closed,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

func NewTickTIckClient() (*TickTickClient, error) {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// 从环境变量获取凭证
	clientID := os.Getenv("TICKTICK_CLIENT_ID")
	clientSecret := os.Getenv("TICKTICK_CLIENT_SECRET")
	accessToken := os.Getenv("TICKTICK_ACCESS_TOKEN")
	refreshToken := os.Getenv("TICKTICK_REFRESH_TOKEN")
	baseURL := os.Getenv("TICKTICK_BASE_URL")
	tokenURL := os.Getenv("TICKTICK_TOKEN_URL")

	// 设置默认值
	if baseURL == "" {
		baseURL = "https://api.dida365.com/open/v1"
	}
	if tokenURL == "" {
		tokenURL = "https://dida365.com/oauth/token"

	}

	// 验证必要凭证
	if accessToken == "" {
		return nil, fmt.Errorf("TICKTICK_ACCESS_TOKEN not set")
	}

	return &TickTickClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		BaseURL:      baseURL,
		TokenURL:     tokenURL,
		HTTPClient:   &http.Client{Timeout: time.Second * 30},
	}, nil
}
func (c *TickTickClient) RefreshAccessToken() error {
	if c.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}
	if c.ClientID == "" || c.ClientSecret == "" {
		return fmt.Errorf("client ID or client secret missing")
	}
	// 准备token请求
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.RefreshToken)

	// 准备Basic Auth验证
	authStr := fmt.Sprintf("%s:%s", c.ClientID, c.ClientSecret)
	authBytes := []byte(authStr)
	authB64 := base64.StdEncoding.EncodeToString(authBytes)

	// 创建请求
	request, err := http.NewRequest("post", c.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Add("Authorization", "Basic "+authB64)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request: %v", response)
	}
	defer response.Body.Close()

	// 检查响应状态
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("error response %s, %s", response.Status, string(body))
	}
	// 解析响应
	var tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	if err = json.NewDecoder(response.Body).Decode(&tokens); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	// 更新令牌
	c.AccessToken = tokens.RefreshToken
	if tokens.RefreshToken != "" {
		c.RefreshToken = tokens.RefreshToken
	}
	// 保存令牌到.env文件中
	return c.saveTokenToEnv()
}

// 保存令牌到.env文件
func (c *TickTickClient) saveTokenToEnv() error {
	// 加载现有的.env文件内容
	envPath := ".env"
	envContent := make(map[string]string)

	// 读取.env文件内容
	if _, err := os.Stat(envPath); err == nil {
		content, err := os.ReadFile(envPath)
		if err != nil {
			return fmt.Errorf("error reading .env file:%v", err)
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				envContent[parts[0]] = parts[1]
			}
		}
	}
	// 更新令牌
	envContent["TICKTICK_ACCESS_TOKEN"] = c.AccessToken
	envContent["TICKTICK_REFRESH_TOKEN"] = c.RefreshToken

	// 确保客户端凭证也被保存
	if c.ClientID != "" {
		envContent["TICKTICK_CLIENT_ID"] = c.ClientID
	}
	if c.ClientSecret != "" {
		envContent["TICKTICK_CLIENT_SECRET"] = c.ClientSecret
	}
	dir := filepath.Dir(envPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}
	// 写入.env文件
	var buf bytes.Buffer
	for key, value := range envContent {
		buf.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	if err := os.WriteFile(envPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing .env file: %v", err)
	}
	return nil
}

func (c *TickTickClient) makeRequest(method, endpoint string, data interface{}) ([]byte, error) {
	// 构建完整URL
	url := c.BaseURL + endpoint

	// 构建请求体
	var reqBody io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshalling data: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// 创建请求

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	// set request header
	request.Header.Set("Authorization", "Bearer "+c.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept-Encoding", "identify")
	request.Header.Set("User-Agent", "ticktick-mcp-go/1.0")
	// send request
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("srror sending request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// 检查是否需要刷新令牌（如果是401未授权）
	if response.StatusCode == http.StatusUnauthorized {
		if err := c.RefreshAccessToken(); err != nil {
			return nil, fmt.Errorf("error refreshing access token: %w", err)
		}

		// 重试请求
		request.Header.Set("Authorization", "Bearer "+c.AccessToken)
		response, err = c.HTTPClient.Do(request)
		if err != nil {
			return nil, fmt.Errorf("error sending request after token refresh: %w", err)
		}
		defer response.Body.Close()

		body, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response after token refresh: %w", err)
		}
	}
	// 检查HTTP状态码
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s, %s", response.Status, string(body))
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
func (c *TickTickClient) GetProjectWithData(projectID string) (map[string]interface{}, error) {
	body, err := c.makeRequest("GET", "/project/"+projectID+"/data", nil)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error unmarshlling project data: %v", err)
	}
	return data, nil
}

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
