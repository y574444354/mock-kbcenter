# 数据库管理

本项目提供了数据库迁移和初始化功能，可以通过命令行工具或Makefile命令来执行。

## 数据库迁移

数据库迁移用于创建或更新数据库表结构，不会添加任何数据。

### 使用Makefile命令

```bash
make db-migrate
```

### 直接使用命令行工具

```bash
go run cmd/dbtools/*.go migrate
```

## 数据库初始化

数据库初始化会先执行迁移，然后添加初始数据（如管理员用户）。

### 使用Makefile命令

```bash
make db-init
```

### 直接使用命令行工具

```bash
go run cmd/dbtools/*.go init
```

## 配置

数据库配置在`config/config.yaml`文件中，可以根据需要修改。系统会自动读取该配置文件，并且如果存在`config/config.local.yaml`文件，会优先使用该文件中的配置进行覆盖。

```yaml
# 数据库配置
database:
  type: postgres  # mysql, postgres, sqlite
  host: localhost
  port: 5432
  user: your_username
  password: your_password
  dbname: go_webserver
  sslmode: disable
  timezone: Asia/Shanghai
```

### 本地配置覆盖

如果需要在本地环境使用不同的配置（例如不同的数据库凭据），可以创建一个`config/config.local.yaml`文件，系统会自动优先使用该文件中的配置。这个文件通常应该被添加到`.gitignore`中，以避免将本地配置提交到版本控制系统。

## 添加新模型

当添加新的数据库模型时，需要在`cmd/dbtools/migrate.go`和`cmd/dbtools/init.go`文件中的`db.AutoMigrate`函数中注册新模型：

```go
// 注册所有需要迁移的模型
if err := db.AutoMigrate(
    &model.User{},
    &model.UserProfile{},
    &model.YourNewModel{}, // 添加新模型
    // 在这里添加其他模型
); err != nil {
    logger.Error("数据库迁移失败", "error", err)
    os.Exit(1)
}