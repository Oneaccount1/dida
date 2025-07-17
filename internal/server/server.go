package server

import (
	"context"
	"dida/internal/client"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var ticktickClient *client.TickTickClient

func InitializeClient() bool {
	var err error
	ticktickClient, err = client.NewTickTIckClient()
	if err != nil {
		log.Printf("Failed to initialize TickTick client: %v", err)
		log.Println("Please run 'ticktick-mcp auth' to authenticate.")
		return false
	}

	// test client connection
	projects, err := ticktickClient.GetProjects()
	if err != nil {
		log.Printf("Failed to access TickTick API: %v", err)
		log.Println("Your access token may have expired. Please run 'ticktick-mcp auth' to refresh it.")
		return false
	}
	log.Printf("Successfully connected to TickTick API with %d projects", len(projects))
	return true
}

func FormatProject(project client.Project) string {
	formatted := fmt.Sprintf("Nmae : %s\n", project.Name)
	formatted += fmt.Sprintf("ID: %s\n", project.ID)

	if project.Color != "" {
		formatted += fmt.Sprintf("Color: %s\n", project.Color)
	}
	if project.ViewMode != "" {
		formatted += fmt.Sprintf("View Mode: %s\n", project.ViewMode)
	}
	formatted += fmt.Sprintf("Closed: %v\n", project.Closed)
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
			return mcp.NewToolResultErrorf(fmt.Sprintf("Error fecting project: %v", err)), nil
		}
		return mcp.NewToolResultText(FormatProject(*project)), nil

	})

	log.Println("Starting TickTick MCP server...")
	return server.ServeStdio(s)
}
