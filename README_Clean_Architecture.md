# TickTick MCP Server - Clean Architecture

这是一个基于 Clean Architecture 设计原则重构的 TickTick MCP 服务器。

## 🏗️ 架构概述

本项目采用 Uncle Bob 的 Clean Architecture 设计，具有以下层次结构：

```
dida/
├── domain/                   # 领域层 - 核心业务规则
│   ├── entities/            # 实体 - 业务对象
│   ├── repositories/        # 仓库接口 - 数据访问契约
│   ├── services/            # 领域服务接口 - 外部服务契约
│   └── errors/              # 业务错误定义
├── usecases/                # 用例层 - 应用业务逻辑
│   ├── task/               # 任务相关用例
│   ├── project/            # 项目相关用例
│   └── auth/               # 认证相关用例
├── adapters/                # 适配器层 - 接口适配
│   ├── repositories/       # 仓库实现
│   └── external/           # 外部服务实现
├── infrastructure/          # 基础设施层 - 框架和工具
│   └── config/             # 配置管理
├── interfaces/              # 接口层 - 外部交互
│   └── mcp/                # MCP接口实现
└── cmd/                     # 应用入口
    └── app/                # 应用程序组装
```

## 🎯 设计原则

### 1. 依赖倒置 (Dependency Inversion)
- 高层模块不依赖低层模块，都依赖抽象
- 抽象不依赖具体实现，具体实现依赖抽象

### 2. 关注点分离 (Separation of Concerns)
- 每层专注于自己的职责
- 业务逻辑与基础设施解耦

### 3. 单一职责 (Single Responsibility)
- 每个类和模块只有一个变化的理由
- 明确的接口定义

### 4. 开放封闭 (Open/Closed)
- 对扩展开放，对修改封闭
- 通过接口实现新功能

## 📦 核心组件

### Domain Layer (领域层)

#### Entities (实体)
```go
// Task - 任务实体
type Task struct {
    ID          string
    ProjectID   string
    Title       string
    // ... 业务属性
}

// 业务方法
func (t *Task) IsCompleted() bool
func (t *Task) IsOverdue() bool
func (t *Task) Complete()
```

#### Repositories (仓库接口)
```go
type TaskRepository interface {
    GetByID(ctx context.Context, taskID string) (*entities.Task, error)
    Create(ctx context.Context, task *entities.Task) error
    // ... 其他方法
}
```

#### Services (服务接口)
```go
type TickTickService interface {
    GetProjects(ctx context.Context) ([]*entities.Project, error)
    CreateTask(ctx context.Context, task *entities.Task) error
    // ... 其他方法
}
```

### Use Cases Layer (用例层)

#### 项目用例
- `GetProjectsUseCase` - 获取项目列表
- `GetProjectUseCase` - 获取单个项目

#### 任务用例
- `GetTasksUseCase` - 获取任务列表
- `CreateTaskUseCase` - 创建任务

#### 认证用例
- `AuthenticateUseCase` - 用户认证
- `RefreshTokenUseCase` - 刷新令牌

### Adapters Layer (适配器层)

#### 仓库实现
- `FileAuthRepository` - 基于文件的认证仓库

#### 外部服务实现
- `TickTickClient` - TickTick API 客户端
- `OAuthAuthService` - OAuth 认证服务

### Infrastructure Layer (基础设施层)

#### 配置管理
- 环境变量加载
- 配置验证
- 默认值设置

### Interfaces Layer (接口层)

#### MCP 接口
- `MCPHandlers` - MCP 工具处理器
- `MCPServer` - MCP 服务器封装

## 🚀 快速开始

### 1. 环境配置

创建 `.env` 文件：
```bash
TICKTICK_CLIENT_ID=your_client_id
TICKTICK_CLIENT_SECRET=your_client_secret
TICKTICK_BASE_URL=https://api.dida365.com/open/v1
TICKTICK_AUTH_URL=https://dida365.com/oauth/authorize
TICKTICK_TOKEN_URL=https://dida365.com/oauth/token
```

### 2. 运行应用

```bash
# 使用新的 Clean Architecture 版本
go run cmd/main_clean.go

# 或者使用原版本进行对比
go run main.go
```

### 3. 使用 MCP 工具

应用启动后，你可以使用以下 MCP 工具：

```bash
# 获取认证 URL
{"method": "tools/call", "params": {"name": "get_auth_url"}}

# 认证
{"method": "tools/call", "params": {"name": "authenticate", "arguments": {"authorization_code": "your_code"}}}

# 获取项目
{"method": "tools/call", "params": {"name": "get_projects"}}

# 获取任务
{"method": "tools/call", "params": {"name": "get_tasks", "arguments": {"project_id": "project_id"}}}

# 创建任务
{"method": "tools/call", "params": {"name": "create_task", "arguments": {"project_id": "project_id", "title": "New Task"}}}
```

## 🧪 测试

### 单元测试
每个用例都包含单元测试，使用 Mock 对象隔离依赖：

```bash
go test ./usecases/...
go test ./domain/...
```

### 集成测试
测试完整的用例流程：

```bash
go test ./cmd/app/...
```

## 📈 与原版本对比

| 特性 | 原版本 | Clean Architecture 版本 |
|------|--------|------------------------|
| **代码组织** | 按技术分层 | 按业务领域分层 |
| **依赖关系** | 紧耦合 | 依赖倒置 |
| **可测试性** | 困难 | 容易（接口 Mock） |
| **可扩展性** | 有限 | 高度可扩展 |
| **业务逻辑位置** | 分散在各层 | 集中在用例层 |
| **错误处理** | 技术错误 | 业务错误 + 技术错误 |

## 🔧 开发指南

### 添加新功能

1. **定义领域实体**（如果需要）
   ```go
   // domain/entities/new_entity.go
   ```

2. **定义仓库接口**
   ```go
   // domain/repositories/new_repository.go
   ```

3. **实现用例**
   ```go
   // usecases/feature/new_usecase.go
   ```

4. **实现适配器**
   ```go
   // adapters/repositories/new_repository_impl.go
   // adapters/external/new_service_impl.go
   ```

5. **添加接口处理器**
   ```go
   // interfaces/mcp/new_handler.go
   ```

6. **在应用中组装**
   ```go
   // cmd/app/wire.go
   ```

### 最佳实践

1. **依赖方向**：始终向内依赖，外层依赖内层
2. **接口隔离**：定义小而专一的接口
3. **错误处理**：使用领域错误，不泄露实现细节
4. **测试优先**：为每个用例编写测试
5. **配置外化**：所有配置通过环境变量管理

## 🎉 优势总结

### 💪 可维护性
- **清晰的职责分离**：每层专注自己的职责
- **低耦合高内聚**：减少模块间依赖
- **易于理解**：业务逻辑集中在用例层

### 🧪 可测试性
- **接口 Mock**：轻松模拟外部依赖
- **独立测试**：每层可单独测试
- **快速反馈**：单元测试运行速度快

### 🚀 可扩展性
- **插拔式架构**：轻松替换实现
- **新功能添加**：不影响现有代码
- **多种接口**：支持 MCP、HTTP、gRPC 等

### 🛡️ 稳定性
- **业务规则保护**：核心逻辑不受外部变化影响
- **错误隔离**：明确的错误边界
- **配置管理**：统一的配置验证和管理

这个重构版本展示了如何将现有的 TickTick MCP 服务器改造为符合 Clean Architecture 原则的高质量代码库。通过明确的分层和依赖管理，我们获得了更好的可维护性、可测试性和可扩展性。