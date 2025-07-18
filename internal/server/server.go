package server

import (
	"context"
	"dida/globalinit"
	"dida/internal/client"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var ticktickClient *client.TickTickClient

func InitializeClient() bool {
	var err error
	logger := globalinit.GetLogger()

	ticktickClient, err = client.NewTickTIckClient()
	if err != nil {
		logger.Errorf("Failed to initialize TickTick client: %v", err)
		logger.Error("Please run 'ticktick-mcp auth' to authenticate.")
		return false
	}

	// test client connection
	projects, err := ticktickClient.GetProjects()
	if err != nil {
		logger.Errorf("Failed to access TickTick API: %v", err)
		logger.Error("Your access token may have expired. Please run 'ticktick-mcp auth' to refresh it.")
		return false
	}
	logger.Infof("Successfully connected to TickTick API with %d projects", len(projects))
	return true
}

func FormatProject(project client.Project) string {
	formatted := fmt.Sprintf("Name: %s\n", project.Name)
	formatted += fmt.Sprintf("ID: %s\n", project.ID)

	if project.Color != "" {
		formatted += fmt.Sprintf("Color: %s\n", project.Color)
	}
	if project.ViewMode != "" {
		formatted += fmt.Sprintf("View Mode: %s\n", project.ViewMode)
	}
	if project.Kind != "" {
		formatted += fmt.Sprintf("Kind: %s\n", project.Kind)
	}
	return formatted
}

// FormatTask 将任务对象格式化为可读字符串
func FormatTask(task client.Task) string {

	formatted := fmt.Sprintf("ID: %s\n", task.ID)
	formatted += fmt.Sprintf("Title: %s\n", task.Title)
	formatted += fmt.Sprintf("Project ID: %s\n", task.ProjectID)

	if task.StartDate != "" {
		formatted += fmt.Sprintf("Start Date: %s\n", task.StartDate)
	}

	if task.DueDate != "" {
		formatted += fmt.Sprintf("Due Date: %s\n", task.DueDate)
	}

	priorityMap := map[int]string{
		0: "None",
		1: "Low",
		3: "Medium",
		5: "High",
	}
	formatted += fmt.Sprintf("Priority: %s\n", priorityMap[task.Priority])

	status := "Active"
	if task.Status == 2 {
		status = "Completed"
	}
	formatted += fmt.Sprintf("Status: %s\n", status)

	if task.Content != "" {
		formatted += fmt.Sprintf("\nContent:\n%s\n", task.Content)
	}

	if len(task.Items) > 0 {
		formatted += fmt.Sprintf("\nSubtasks (%d):\n", len(task.Items))
		for i, item := range task.Items {
			statusMark := "□"
			if item.Status == 1 {
				statusMark = "✓"
			}
			formatted += fmt.Sprintf("%d. [%s] %s\n", i+1, statusMark, item.Title)
		}
	}

	return formatted
}

func Start() error {
	// 获取日志器
	logger := globalinit.GetLogger()

	// 初始化TickTick客户端
	if !InitializeClient() {
		return fmt.Errorf("fail to initialize TickTick client")
	}
	// 创建MCP服务器
	s := server.NewMCPServer(
		"TickTick MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// 添加工具：获取所有项目
	getProjectsTool := mcp.NewTool("get_projects",
		mcp.WithDescription("Get all projects from TickTick."),
	)
	s.AddTool(getProjectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Fail to initialize TickTick client.Please check your API credentials."), nil
			}
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
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
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
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
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
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
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
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
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
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
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
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		tasksID, err := request.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		if err := ticktickClient.CompletedTask(projectID, tasksID); err != nil {
			return mcp.NewToolResultErrorf("Fail to complete task: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Task completed successfully!\n")), nil
	})

	// 删除任务
	deleteTaskTool := mcp.NewTool("delete_task",
		mcp.WithDescription("Delete a task"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("the project_id containing this task"),
		),
		mcp.WithString("task_id",
			mcp.Required(),
			mcp.Description("the iD of task"),
		),
	)

	s.AddTool(deleteTaskTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if ticktickClient == nil {
			if !InitializeClient() {
				return mcp.NewToolResultError("Failed to initialize TickTick client. Please check your API credentials."), nil
			}
		}
		projectID, err := request.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}
		task_id, err := request.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultErrorf(err.Error()), nil
		}

		if err := ticktickClient.DeleteTask(projectID, task_id); err != nil {
			return mcp.NewToolResultErrorf("Fail to delete task: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Task deleted successfully!\n")), nil
	})

	logger.Info("Starting TickTick MCP server...")

	return server.ServeStdio(s)
}
