# TickTick MCP Server

一个为 TickTick (滴答清单) 提供 MCP (Model Context Protocol) 接口的 Go 服务器，让 AI 助手能够直接管理您的 TickTick 任务和项目。

## 功能特性

- 🔐 **完整的 OAuth2 认证流程** - 安全的用户授权和令牌管理
- 📋 **项目和任务的完整 CRUD 操作** - 创建、读取、更新、删除任务和项目
- 🔄 **智能令牌刷新机制** - 自动处理访问令牌过期和刷新
- 📝 **结构化日志记录** - 详细的操作日志和错误追踪
- 🛡️ **统一的错误处理** - 一致的错误响应和用户友好的错误信息
- ⚙️ **灵活的配置管理** - 环境变量配置，内置合理默认值
- 🧹 **Clean Architecture 设计** - 模块化架构，易于维护和扩展
- 🔧 **MCP Inspector 支持** - 内置调试和测试工具支持

## 支持的 MCP 工具

| 工具名称 | 描述 | 参数 |
|---------|------|------|
| `oauth_authorize` | 启动 OAuth2 授权流程 | 无 |
| `get_projects` | 获取所有项目 | 无 |
| `get_project` | 获取特定项目详情 | `project_id` |
| `get_project_tasks` | 获取项目中的所有任务 | `project_id` |
| `get_task` | 获取特定任务详情 | `project_id`, `task_id` |
| `create_task` | 创建新任务 | `project_id`, `title`, `content?`, `start_date?`, `due_date?`, `priority?` |
| `update_task` | 更新任务 | `task_id`, `project_id`, `title`, `content?`, `start_date?`, `due_date?`, `priority?` |
| `complete_task` | 完成任务 | `project_id`, `task_id` |
| `delete_task` | 删除任务 | `project_id`, `task_id` |

## 快速开始

### 1. 前置要求

- Go 1.23.4 或更高版本
- TickTick 开发者账号和 API 凭证

### 2. 获取 TickTick API 凭证

1. 访问 [TickTick 开发者中心](https://developer.ticktick.com)
2. 创建新的应用程序
3. 获取 `Client ID` 和 `Client Secret`
4. 设置回调 URL 为 `http://localhost:8000/callback`

### 3. 安装和配置

```bash
# 克隆项目
git clone <repository-url>
cd dida

# 构建项目
go build -o dida.exe ./cmd/ticktick-mcp
```

### 4. 环境变量配置

创建 `.env` 文件（只需配置认证信息）：

```env
# TickTick API 认证配置
TICKTICK_CLIENT_ID=your_client_id
TICKTICK_CLIENT_SECRET=your_client_secret

# 认证后自动生成（无需手动设置）
TICKTICK_ACCESS_TOKEN=
TICKTICK_REFRESH_TOKEN=
```

**重要说明**:
- ✅ **只需配置** `CLIENT_ID` 和 `CLIENT_SECRET`
- ✅ **访问令牌自动获取** - 通过 OAuth2 流程自动生成和保存
- ✅ **内置默认配置** - API 端点、服务器设置等已预配置
- ✅ **安全存储** - 所有敏感信息仅存储在本地 `.env` 文件中

## 使用方法

### 启动服务器

```bash
# 直接启动服务器
./dida.exe

# 或使用开发模式
go run ./cmd/ticktick-mcp
```

### OAuth2 授权流程

1. **配置环境变量**: 确保 `.env` 文件中已配置 `TICKTICK_CLIENT_ID` 和 `TICKTICK_CLIENT_SECRET`

2. **启动授权**: 在 AI 助手中调用 `oauth_authorize` 工具

3. **完成授权**:
   - AI 助手会提供一个授权 URL
   - 访问该 URL 并登录您的 TickTick 账号
   - 授权完成后，访问令牌会自动保存到 `.env` 文件

4. **开始使用**: 授权完成后即可使用其他 MCP 工具管理任务

### 调试和测试

使用 MCP Inspector 进行调试：

```bash
# 启动 Inspector 调试界面
npx @modelcontextprotocol/inspector go run ./cmd/ticktick-mcp

# 或使用已构建的可执行文件
npx @modelcontextprotocol/inspector ./dida.exe
```

Inspector 提供：
- 🔧 **工具测试界面** - 可视化测试所有 MCP 工具
- 📊 **实时日志** - 查看请求/响应和错误信息
- 🔍 **参数验证** - 验证工具参数格式
- 📈 **性能监控** - 监控 API 调用性能

### 与 AI 助手集成

服务器启动后，它会通过标准输入/输出与支持 MCP 协议的 AI 助手通信。确保您的 AI 助手配置正确指向此服务器。

## 项目结构

```
dida/
├── cmd/ticktick-mcp/          # 应用程序入口
├── internal/                  # 内部包
│   ├── auth/                  # OAuth2 认证管理
│   │   └── auth.go           # 认证流程、令牌刷新
│   ├── client/                # TickTick API 客户端
│   │   ├── ticktick_client.go # API 客户端实现
│   │   └── model.go          # 数据模型定义
│   ├── config/                # 配置管理
│   │   └── config.go         # 配置加载和验证
│   ├── errors/                # 错误处理
│   │   └── errors.go         # 统一错误定义
│   ├── logger/                # 日志记录
│   │   └── logger.go         # 结构化日志实现
│   └── server/                # MCP 服务器（重构后）
│       ├── server.go         # 服务器核心逻辑
│       ├── tools.go          # MCP 工具定义
│       └── help.go           # 辅助格式化函数
├── globalinit/                # 全局初始化
│   └── init.go               # 全局组件初始化
├── go.mod                     # Go 模块定义
├── go.sum                     # 依赖校验和
├── README.md                  # 项目文档
├── .env.example               # 环境变量示例
└── .gitignore                 # Git 忽略文件
```

### 架构特点

- 🏗️ **Clean Architecture** - 清晰的分层架构，依赖关系明确
- 📦 **模块化设计** - 每个包都有单一职责
- 🔄 **依赖注入** - 松耦合的组件设计
- 🧪 **易于测试** - 接口驱动的设计便于单元测试

 
## 故障排除

### 常见问题

#### 1. **认证相关问题**

**问题**: `Client ID or Client Secret not found`
```bash
解决方案:
1. 检查 .env 文件是否存在
2. 确认 TICKTICK_CLIENT_ID 和 TICKTICK_CLIENT_SECRET 已正确配置
3. 重启服务器
```

**问题**: `Authorization failed` 或回调超时
```bash
解决方案:
1. 确认 TickTick 应用回调 URL 设置为: http://localhost:8000/callback
2. 检查防火墙是否阻止了 8000 端口
3. 确保浏览器能正常访问 localhost:8000
```

#### 2. **令牌相关问题**

**问题**: `Access token expired and no refresh token available`
```bash
解决方案:
1. 使用 oauth_authorize 工具重新授权
2. 确保授权流程完整完成
3. 检查 .env 文件中的 REFRESH_TOKEN 是否已保存
```

**问题**: `Failed to refresh access token`
```bash
解决方案:
1. 重新运行 OAuth2 授权流程
2. 检查网络连接
3. 确认 Client Secret 未过期
```

#### 3. **API 请求问题**

**问题**: `Error fetching projects` 或其他 API 错误
```bash
解决方案:
1. 检查网络连接
2. 确认 TickTick API 服务状态
3. 查看 log.txt 文件获取详细错误信息
4. 尝试重新授权
```

#### 4. **Inspector 调试问题**

**问题**: Inspector 无法启动
```bash
解决方案:
1. 确保 Node.js 已安装: node --version
2. 清除 npm 缓存: npm cache clean --force
3. 重新安装: npm install -g @modelcontextprotocol/inspector
```

### 日志查看

```bash
# 查看实时日志
tail -f log.txt

# 查看错误日志
grep "ERROR" log.txt
```

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本仓库
2. 创建功能分支: `git checkout -b feature/your-feature`
3. 提交更改: `git commit -am 'Add some feature'`
4. 推送分支: `git push origin feature/your-feature`
5. 提交 Pull Request

### 代码规范

- 遵循 Go 语言官方代码规范
- 使用 `go fmt ./...` 格式化代码
- 使用 `go vet ./...` 进行代码检查
- 为新功能添加相应的测试

## 许可证

[MIT License](LICENSE)
