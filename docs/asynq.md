# Asynq 异步任务队列使用指南

Asynq 是一个基于 Redis 的 Go 异步任务队列库，本项目已集成 Asynq 用于处理后台任务。

## 安装和配置

Asynq 已作为项目依赖自动安装，配置位于 `config/config.yaml` 文件中：

```yaml
asynq:
  concurrency: 10      # 并发worker数量
  retry_count: 3       # 最大重试次数
  retry_delay: 10      # 重试延迟(秒)
  redis_pool_size: 20  # Redis连接池大小
  queues:              # 队列优先级
    default: 5
    critical: 10
  log:                 # Asynq日志配置
    level: info
    format: json
    output_path: ./logs/asynq.log
    max_size: 100
    max_backups: 10
    max_age: 30
    compress: true
```

## 任务定义

### 1. 定义任务类型

在 `tasks/types.go` 中定义任务类型常量：

```go
const (
    TypeEmailDelivery = "email:deliver"
    TypeImageResize   = "image:resize"
)
```

### 2. 定义任务负载结构

每个任务需要定义自己的负载结构，例如邮件任务：

```go
type EmailDeliveryPayload struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}
```

### 3. 创建任务函数

为每个任务类型创建任务创建函数：

```go
func NewEmailDeliveryTask(payload EmailDeliveryPayload, queue string) (*asynq.Task, error) {
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("序列化邮件任务负载失败: %w", err)
    }
    return asynq.NewTask(TypeEmailDelivery, payloadBytes, asynq.Queue(queue)), nil
}
```

### 4. 任务处理器

实现任务处理逻辑：

```go
func HandleEmailDeliveryTask(ctx context.Context, task *asynq.Task) error {
    var payload EmailDeliveryPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("反序列化邮件任务负载失败: %w", err)
    }
    
    // 实现实际的任务处理逻辑
    return nil
}
```

## 注册任务处理器

在 worker 启动时注册任务处理器：

```go
mux := asynq.NewServeMux()
mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
mux.HandleFunc(tasks.TypeImageResize, tasks.HandleImageResizeTask)
```

## 客户端API使用

### 1. 创建任务

```go
payload := tasks.EmailDeliveryPayload{
    To:      "user@example.com",
    Subject: "欢迎邮件",
    Body:    "感谢注册我们的服务",
}

task, err := tasks.NewEmailDeliveryTask(payload, "critical")
if err != nil {
    // 处理错误
}
```

### 2. 添加任务到队列

```go
client := asynq.NewClient(asynq.RedisClientOpt{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
defer client.Close()

info, err := client.Enqueue(task)
if err != nil {
    // 处理错误
}
```

## 启动Worker

使用 Makefile 命令启动 worker：

```bash
make worker
```

或直接运行：

```bash
go run cmd/worker/main.go
```

## 监控和管理

Asynq 提供了 Web UI 用于监控和管理任务队列：

```go
mux := http.NewServeMux()
asynqmon := asynqmon.New(asynqmon.Options{
    RootPath:     "/monitoring", // 根路径
    RedisConnOpt: asynq.RedisClientOpt{Addr: "localhost:6379"},
})
mux.Handle("/monitoring/", http.StripPrefix("/monitoring", asynqmon))
```

访问 `http://localhost:8080/monitoring` 查看任务队列状态。

## 示例代码

### 发送邮件任务

```go
payload := tasks.EmailDeliveryPayload{
    To:      "user@example.com",
    Subject: "测试邮件",
    Body:    "这是一封测试邮件",
}

task, err := tasks.NewEmailDeliveryTask(payload, "critical")
if err != nil {
    log.Fatal(err)
}

client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
defer client.Close()

if _, err := client.Enqueue(task); err != nil {
    log.Fatal(err)
}
```

### 图片处理任务

```go
payload := tasks.ImageResizePayload{
    SourceURL: "https://example.com/image.jpg",
    Width:     800,
    Height:    600,
}

task, err := tasks.NewImageResizeTask(payload, "default")
if err != nil {
    log.Fatal(err)
}

client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
defer client.Close()

if _, err := client.Enqueue(task); err != nil {
    log.Fatal(err)
}