# 构建阶段
FROM golang:1.24-alpine AS builder

# 定义构建参数
ARG GOPROXY=https://goproxy.cn,direct
ARG ALPINE_MIRROR=https://mirrors.aliyun.com
ARG BUILD_DBTOOLS=true
ARG BUILD_REDISTOOLS=true

WORKDIR /app

# 设置构建环境
RUN sed -i "s|https://dl-cdn.alpinelinux.org|${ALPINE_MIRROR}|g" /etc/apk/repositories && \
    apk add --no-cache gcc musl-dev && \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=${GOPROXY} && \
    go env -w GOSUMDB=off

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并构建
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-webserver . && \
    mkdir -p bin && \
    if [ "$BUILD_DBTOOLS" = "true" ]; then \
        go build -o bin/dbtools ./cmd/dbtools; \
    fi && \
    if [ "$BUILD_REDISTOOLS" = "true" ]; then \
        go build -o bin/redistools ./cmd/redistools; \
    fi

# 运行阶段
FROM alpine:latest

# 定义构建参数
ARG ALPINE_MIRROR=https://mirrors.aliyun.com

# 设置运行时环境
RUN sed -i "s|https://dl-cdn.alpinelinux.org|${ALPINE_MIRROR}|g" /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser && \
    mkdir -p /app/{logs,config,bin} && \
    chown -R appuser:appuser /app

# 设置时区和用户
ENV TZ=Asia/Shanghai
USER appuser
WORKDIR /app
# 复制构建产物
COPY --from=builder --chown=appuser:appuser /app/go-webserver .
COPY --from=builder --chown=appuser:appuser /app/bin ./bin


# 暴露端口和健康检查
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O- http://localhost:8080/health || exit 1

CMD ["./go-webserver"]