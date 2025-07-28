package mcp

import (
	"log"
	
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer MCP服务器
type MCPServer struct {
	server   *server.MCPServer
	handlers *MCPHandlers
}

// NewMCPServer 创建MCP服务器
func NewMCPServer(name, version string, handlers *MCPHandlers) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	
	mcpServer := &MCPServer{
		server:   s,
		handlers: handlers,
	}
	
	mcpServer.registerTools()
	
	return mcpServer
}

// registerTools 注册工具
func (s *MCPServer) registerTools() {
	// 获取项目列表工具
	getProjectsTool := mcp.NewTool("get_projects",
		mcp.WithDescription("Get all projects from TickTick."),
		mcp.WithString("include_inactive",
			mcp.Description("Include inactive/closed projects (default: false)"),
		),
	)
	s.server.AddTool(getProjectsTool, s.handlers.GetProjects)
	
	// 获取单个项目工具
	getProjectTool := mcp.NewTool("get_project",
		mcp.WithDescription("Get details about a specific project."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project"),
		),
	)
	s.server.AddTool(getProjectTool, s.handlers.GetProject)
	
	// 获取任务列表工具
	getTasksTool := mcp.NewTool("get_tasks",
		mcp.WithDescription("Get tasks from TickTick."),
		mcp.WithString("project_id",
			mcp.Description("ID of the project (optional)"),
		),
		mcp.WithString("include_completed",
			mcp.Description("Include completed tasks (default: false)"),
		),
		mcp.WithString("priority",
			mcp.Description("Filter by priority (0=None, 1=Low, 2=Medium, 3=High)"),
		),
	)
	s.server.AddTool(getTasksTool, s.handlers.GetTasks)
	
	// 创建任务工具
	createTaskTool := mcp.NewTool("create_task",
		mcp.WithDescription("Create a new task in TickTick."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("ID of the project"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the task"),
		),
		mcp.WithString("content",
			mcp.Description("Content/description of the task"),
		),
		mcp.WithString("description",
			mcp.Description("Additional description"),
		),
		mcp.WithString("priority",
			mcp.Description("Priority (0=None, 1=Low, 2=Medium, 3=High)"),
		),
	)
	s.server.AddTool(createTaskTool, s.handlers.CreateTask)
	
	// 获取认证URL工具
	getAuthURLTool := mcp.NewTool("get_auth_url",
		mcp.WithDescription("Get the authentication URL for TickTick OAuth."),
	)
	s.server.AddTool(getAuthURLTool, s.handlers.GetAuthURL)
	
	// 认证工具
	authenticateTool := mcp.NewTool("authenticate",
		mcp.WithDescription("Authenticate with TickTick using authorization code."),
		mcp.WithString("authorization_code",
			mcp.Required(),
			mcp.Description("Authorization code from OAuth callback"),
		),
	)
	s.server.AddTool(authenticateTool, s.handlers.Authenticate)
	
	log.Println("Registered MCP tools:")
	log.Println("- get_projects: Get all projects")
	log.Println("- get_project: Get specific project details")
	log.Println("- get_tasks: Get tasks (optionally filtered)")
	log.Println("- create_task: Create a new task")
	log.Println("- get_auth_url: Get OAuth authentication URL")
	log.Println("- authenticate: Authenticate with authorization code")
}

// Start 启动服务器
func (s *MCPServer) Start() error {
	log.Println("Starting TickTick MCP server...")
	return server.ServeStdio(s.server)
}