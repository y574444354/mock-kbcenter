# Go WebServer

基于Gin框架的Go Web服务器脚手架，具有完整的分层设计和丰富的功能。

## 功能特性

- 使用 [Gin](https://github.com/gin-gonic/gin) 框架构建高性能Web服务
- 支持国际化 (i18n)，使用 [go-i18n](https://github.com/nicksnyder/go-i18n)
- 集成 [GORM](https://gorm.io/) 进行数据库操作，支持MySQL、PostgreSQL和SQLite
- 支持 [Redis](https://github.com/go-redis/redis) 缓存
- 使用 [Zap](https://github.com/uber-go/zap) 进行高性能日志记录
- 使用 [YAML](https://github.com/go-yaml/yaml) 进行配置管理
- 集成 [Swagger](https://github.com/swaggo/gin-swagger) 自动生成API文档
- 支持作为Windows服务运行
- 支持Docker容器化部署
- 使用 [Asynq](https://github.com/hibiken/asynq) 处理异步任务，支持独立 worker 进程
- 符合SOLID设计原则的清晰分层架构

## 项目结构

```
.
├── api                 # API层，处理HTTP请求
│   └── v1              # API版本1
├── config              # 配置文件和配置管理
├── docs                # Swagger文档
├── i18n                # 国际化资源
│   └── locales         # 语言文件
├── internal            # 内部包
│   ├── middleware      # HTTP中间件
│   ├── model           # 数据模型
│   ├── repository      # 数据访问层
│   └── service         # 业务逻辑层
├── pkg                 # 通用包
│   ├── db              # 数据库连接管理
│   ├── logger          # 日志管理
│   ├── redis           # Redis连接管理
│   ├── asynq           # 异步任务处理
│   └── utils           # 工具函数
└── scripts             # 脚本文件
```

## 分层设计

项目采用清晰的分层设计，遵循SOLID原则：

1. **交互层 (API)**: 处理HTTP请求和响应，参数验证，路由管理
2. **业务逻辑层 (Service)**: 实现业务逻辑，协调各种资源
3. **数据访问层 (Repository)**: 处理数据存储和检索，封装数据库操作
4. **模型层 (Model)**: 定义数据结构和业务实体
5. **通用机制层 (pkg)**: 提供通用功能，如日志、数据库连接、工具函数等

## 开发规范

本项目采用主干开发模式(Trunk-based Development)，详细的分支管理和提交规范请参阅[分支管理规范](./docs/branch_guidelines.md)。

## 快速开始

### 前置条件

- Go 1.24 或更高版本
- MySQL/PostgreSQL/SQLite (根据配置选择)
- Redis (可选)

### 安装

1. 克隆仓库

```bash
git clone https://github.com/yourusername/mock-kbcenter.git
cd mock-kbcenter
```

2. 安装依赖

```bash
go mod download
```

3. 修改配置文件

```bash
cp config/config.yaml config/config.local.yaml
# 编辑 config/config.local.yaml 设置你的配置
```

4. 运行服务器

```bash
go run main.go --config ./config/config.local.yaml
```

### 异步任务处理

项目使用 Asynq 处理异步任务，如邮件发送、图片处理等。worker 进程可以独立运行。

#### 运行 worker

```bash
# 构建 worker
make build-worker

# 运行 worker
make run-worker
```

#### 使用 Docker 运行 worker

```bash
# 构建 worker Docker 镜像
make docker-build-worker

# 运行 worker 容器
docker run -d --name mock-kbcenter-worker mock-kbcenter-worker
```

### 使用Docker

1. 构建Docker镜像

```bash
docker build -t mock-kbcenter .
```

2. 运行容器

```bash
docker run -p 8080:8080 mock-kbcenter
```

### 作为Windows服务运行

1. 安装服务

```bash
go run scripts/windows_service.go -service=install
```

2. 启动服务

```bash
go run scripts/windows_service.go -service=start
```

3. 停止服务

```bash
go run scripts/windows_service.go -service=stop
```

4. 卸载服务

```bash
go run scripts/windows_service.go -service=uninstall
```

## API文档

启动服务后，访问 http://localhost:8080/swagger/index.html 查看API文档。

生成Swagger文档：

```bash
# 安装swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g main.go -o docs
```

## 国际化

项目支持多语言，语言文件位于 `i18n/locales` 目录。

添加新语言：

1. 创建新的语言文件，如 `i18n/locales/fr.yaml`
2. 在代码中使用翻译函数：

```go
// 在中间件中设置语言
c.Set("locale", "fr")

// 在处理器中使用翻译
translate := c.MustGet("translate").(func(string, map[string]interface{}) string)
message := translate("welcome.message", nil)
```

## 许可证

[MIT](LICENSE)