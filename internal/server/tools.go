package server

import (
	"context"
	"dida/globalinit"
	"dida/internal/auth"
	"dida/internal/client"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"os"
)

func InitAllTools(s *server.MCPServer) error {
	// 获取日志器
	logger := globalinit.GetLogger()
	// 添加工具：获取所有项目
	getProjectsTool := mcp.NewTool("get_projects",
		mcp.WithDescription("Get all projects from TickTick."),
	)
	s.AddTool(getProjectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		projects, err := ticktickClient.GetProjects()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error fetching projects: %v", err)), nil
		}

		if len(projects) == 0 {
			return mcp.NewToolResultText("No projects found."), nil
		}

		result := fmt.Sprintf("Found %d projects:\n\n", len(projects))

		for i, project := range projects {
			result += fmt.Sprintf("Project %d:\n%s\n", i+1, FormatProject(project))
		}
		return mcp.NewToolResultText(result), nil
	})

	// 获取特定项目
	getProjectTool := mcp.NewTool("get_project",
		mcp.WithDescription("Get details about a specific project."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project"),
		),
	)
	s.AddTool(getProjectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// 获取项目ID
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		// 获取项目
		project, err := ticktickClient.GetProject(projectID)
		if err != nil {
			return mcp.NewToolResultErrorf(fmt.Sprintf("Error fetching project: %v", err)), nil
		}
		return mcp.NewToolResultText(FormatProject(*project)), nil

	})

	// 获取所有任务在指定Project中
	getProjectTasks := mcp.NewTool("get_project_tasks",
		mcp.WithDescription("Get all tasks from a specific project"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project to retrieve tasks from"),
		),
	)
	s.AddTool(getProjectTasks, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		// 获取projectID
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		// 获取任务
		projectData, err := ticktickClient.GetProjectWithData(projectID)
		if err != nil {
			return mcp.NewToolResultErrorf(fmt.Sprintf("Error fetching project data: %v", err)), nil
		}
		result := ""
		if len(projectData.Tasks) == 0 {
			result = "No tasks found in project"
		} else {
			result = fmt.Sprintf("Found %d tasks:\n\n", len(projectData.Tasks))
			for i, task := range projectData.Tasks {
				result += fmt.Sprintf("Task %d: \n%s\n", i+1, FormatTask(task))
			}
		}
		return mcp.NewToolResultText(result), nil
	})

	// 获取指定Project的指定Task
	getTask := mcp.NewTool("get_task",
		mcp.WithDescription("Get details about a specific task"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project"),
		),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("ID of the task"),
		),
	)
	s.AddTool(getTask, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		taskID, err := request.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		task, err := ticktickClient.GetTask(projectID, taskID)
		if err != nil {
			return mcp.NewToolResultErrorf("Error fetching task: %v", err), nil
		}

		return mcp.NewToolResultText(FormatTask(*task)), nil
	})
	// 创建任务
	createTaskTool := mcp.NewTool("create_task",
		mcp.WithDescription("Create a new task in a specific project"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project to add the task to"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the task"),
		),
		mcp.WithString("content",
			mcp.Description("Content/description of the task"),
		),
		mcp.WithString("start_date",
			mcp.Description("Start date in format YYYY-MM-DDThh:mm:ssZ"),
		),
		mcp.WithString("due_date",
			mcp.Description("Due date in format YYYY-MM-DDThh:mm:ssZ"),
		),
		mcp.WithString("priority",
			mcp.Description("Priority level: 0=None, 1=Low, 3=Medium, 5=High"),
		),
	)
	s.AddTool(createTaskTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		task := client.Task{
			ProjectID: projectID,
			Title:     title,
		}

		// 获取可选选项
		content := request.GetString("content", "")
		if content != "" {
			task.Content = content
		}

		startDate := request.GetString("start_date", "")
		if startDate != "" {
			task.StartDate = startDate
		}
		dueDate := request.GetString("due_date", "")
		if dueDate != "" {
			task.DueDate = dueDate
		}
		priority := request.GetInt("priority", 0)
		task.Priority = priority
		createdTask, err := ticktickClient.CreateTask(task)
		if err != nil {
			return mcp.NewToolResultErrorf("Failed to create task: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Task created successfully:\n%s", FormatTask(*createdTask))), nil

	})

	// 更新任务
	updateTaskTool := mcp.NewTool("update_task",
		mcp.WithDescription("Update an existing task"),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("ID of the task to update"),
		),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project containing the task"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the task"),
		),
		mcp.WithString("content",
			mcp.Description("Content/description of the task"),
		),
		mcp.WithString("start_date",
			mcp.Description("Start date in format YYYY-MM-DDThh:mm:ssZ"),
		),
		mcp.WithString("due_date",
			mcp.Description("Due date in format YYYY-MM-DDThh:mm:ssZ"),
		),
		mcp.WithString("priority",
			mcp.Description("Priority level: 0=None, 1=Low, 3=Medium, 5=High"),
		),
	)
	s.AddTool(updateTaskTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		// 获取请求参数
		taskID, err := request.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		task := client.Task{
			ID:        taskID,
			ProjectID: projectID,
		}

		// 获取可选选项
		content := request.GetString("content", "")
		if content != "" {
			task.Content = content
		}

		startDate := request.GetString("start_date", "")
		if startDate != "" {
			task.StartDate = startDate
		}
		dueDate := request.GetString("due_date", "")
		if dueDate != "" {
			task.DueDate = dueDate
		}
		priority := request.GetInt("priority", 0)
		task.Priority = priority

		updatedTask, err := ticktickClient.UpdateTask(task)
		if err != nil {
			return mcp.NewToolResultErrorf("Failed to update task: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Task updated successfully:\n%s", FormatTask(*updatedTask))), nil
	})

	// 完成任务
	completeTaskTool := mcp.NewTool("complete_task",
		mcp.WithDescription("Mark a task as completed"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project containing the task"),
		),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("ID of the task to mark as completed"),
		),
	)
	s.AddTool(completeTaskTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		taskID, err := request.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		if err := ticktickClient.CompletedTask(projectID, taskID); err != nil {
			return mcp.NewToolResultErrorf("Failed to complete task: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Task completed successfully!\n")), nil
	})

	// 删除任务
	deleteTaskTool := mcp.NewTool("delete_task",
		mcp.WithDescription("Delete a task"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project containing this task"),
		),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("ID of the task to delete"),
		),
	)

	s.AddTool(deleteTaskTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := ensureClientInitialized(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		taskID, err := request.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		if err := ticktickClient.DeleteTask(projectID, taskID); err != nil {
			return mcp.NewToolResultErrorf("Failed to delete task: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Task deleted successfully!\n")), nil
	})

	// 添加OAuth2授权工具
	oauthTool := mcp.NewTool("oauth_authorize",
		mcp.WithDescription("Start OAuth2 authorization flow for TickTick. This will provide a URL for the user to visit and complete authorization."),
	)
	s.AddTool(oauthTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 从环境变量读取Client ID和Secret
		clientID := os.Getenv("TICKTICK_CLIENT_ID")
		clientSecret := os.Getenv("TICKTICK_CLIENT_SECRET")

		if clientID == "" || clientSecret == "" {
			return mcp.NewToolResultError("Client ID or Client Secret not found in environment variables. Please check your .env file."), nil
		}

		// 创建认证管理器
		tickAuth, err := auth.NewTickTickAuth(clientID, clientSecret)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to initialize authentication: %v", err)), nil
		}

		// 生成授权URL
		authURL := tickAuth.GetAuthURL()

		result := fmt.Sprintf(`🔐 TickTick OAuth2 Authorization Required

Please visit the following URL to authorize this application:

%s

Instructions:
1. Click the URL above or copy and paste it into your browser
2. Log in to your TickTick account if prompted
3. Grant the requested permissions
4. You will be redirected to a callback page
5. The authorization will be completed automatically

Note: Make sure your TickTick application's callback URL is set to: http://localhost:8000/callback

Waiting for authorization...`, authURL)

		// 启动认证流程（这会启动本地服务器等待回调）
		go func() {
			if err := tickAuth.StartAuthFlow(); err != nil {
				logger.Errorf("OAuth2 authorization failed: %v", err)
			} else {
				logger.Info("OAuth2 authorization completed successfully")
			}
		}()

		return mcp.NewToolResultText(result), nil
	})

	return nil
}
