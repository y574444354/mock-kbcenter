# Go Web服务器项目Makefile

# 检查.env文件是否存在，如果不存在则提示
ifeq (,$(wildcard .env))
$(info 注意：.env文件不存在，将使用默认值。您可以复制.env.example创建自己的.env文件。)
endif

# 加载.env文件（如果存在）
-include .env

# 变量定义（.env文件中的变量会覆盖这些默认值）
APP_NAME ?= go-webserver
MAIN_FILE := main.go
DOCKER_IMAGE ?= $(APP_NAME):latest
DOCKER_CONTAINER := $(APP_NAME)
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*")
GOPATH := $(shell go env GOPATH)
# 默认GOPROXY设置（可通过.env文件或命令行参数覆盖）
GOPROXY ?= https://goproxy.cn,direct
# 默认Alpine镜像源设置
ALPINE_MIRROR ?= https://mirrors.aliyun.com

# 默认目标
.PHONY: all
all: build build-worker

# 帮助信息
.PHONY: env help
env:
	@echo "当前环境变量设置："
	@echo "APP_NAME: $(APP_NAME)"
	@echo "DOCKER_IMAGE: $(DOCKER_IMAGE)"
	@echo "GOPROXY: $(GOPROXY)"
	@echo "ALPINE_MIRROR: $(ALPINE_MIRROR)"

.PHONY: help
help:
	@echo "Go Web服务器项目管理命令："
	@echo "make build         - 构建主应用程序"
	@echo "make run           - 运行主应用程序"
	@echo "make run-worker    - 运行worker进程"
	@echo "make test          - 执行测试"
	@echo "make clean         - 清理生成的文件"
	@echo "make fmt           - 格式化代码"
	@echo "make lint          - 检查代码质量"
	@echo "make swagger       - 生成Swagger文档"
	@echo "make db-migrate    - 执行数据库迁移"
	@echo "make db-init       - 初始化数据库"
	@echo "make redis-clear   - 清除Redis缓存"
	@echo "make docker-build  - 构建Docker镜像"
	@echo "make env           - 显示当前环境变量设置"

# 构建主应用程序
.PHONY: build
build:
	@echo "构建主应用程序..."
	@go build -o $(APP_NAME) $(MAIN_FILE)
	@echo "主应用程序构建完成: $(APP_NAME)"

# 运行主应用程序
.PHONY: run
run:
	@echo "运行主应用程序..."
	@go run $(MAIN_FILE) web

# 运行worker进程
.PHONY: run-worker
run-worker:
	@echo "运行worker进程..."
	@go run $(MAIN_FILE) worker

# 执行测试
.PHONY: test
test:
	@echo "执行测试..."
	@go test -v ./...

# 清理生成的文件
.PHONY: clean
clean:
	@echo "清理生成的文件..."
	@rm -f $(APP_NAME) worker
	@go clean
	@echo "清理完成"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	@gofmt -s -w $(GO_FILES)
	@echo "格式化完成"

# 检查代码质量
.PHONY: lint
lint:
	@echo "检查代码质量..."
	@if [ ! -f $(GOPATH)/bin/golangci-lint ]; then \
		echo "安装 golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ./...

# 生成Swagger文档
.PHONY: swagger
swagger:
	@echo "生成Swagger文档..."
	@if [ ! -f $(GOPATH)/bin/swag ]; then \
		echo "安装 swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@swag init -g main.go -o swagger
	@echo "Swagger文档生成完成"

# 数据库相关命令
.PHONY: db-migrate
db-migrate:
	@echo "执行数据库迁移..."
	@go run cmd/dbtools/*.go migrate
	@echo "数据库迁移完成"

.PHONY: db-init
db-init:
	@echo "初始化数据库..."
	@go run cmd/dbtools/*.go init
	@echo "数据库初始化完成"

# Docker相关命令
.PHONY: docker-build
docker-build:
	@echo "构建Docker镜像..."
	@echo "使用GOPROXY: $(GOPROXY)"
	@echo "使用ALPINE_MIRROR: $(ALPINE_MIRROR)"
	@docker build --build-arg GOPROXY=$(GOPROXY) --build-arg ALPINE_MIRROR=$(ALPINE_MIRROR) -t $(DOCKER_IMAGE) -f docker/Dockerfile .
	@echo "Docker镜像构建完成: $(DOCKER_IMAGE)"

# Redis相关命令
.PHONY: redis-clear
redis-clear:
	@echo "清除Redis缓存..."
	@go run cmd/redistools/*.go
	@echo "Redis缓存清除完成"
