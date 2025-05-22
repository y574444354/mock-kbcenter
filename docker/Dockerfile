FROM golang:1.24-alpine AS builder

# 定义构建参数
ARG GOPROXY=https://goproxy.cn,direct
ARG ALPINE_MIRROR=mirrors.aliyun.com

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN sed -i "s#dl-cdn.alpinelinux.org#${ALPINE_MIRROR}#g" /etc/apk/repositories && \
   apk update && apk add --no-cache gcc musl-dev

# 复制go.mod和go.sum
COPY go.mod ./
# COPY go.sum ./
RUN go env -w GO111MODULE=on && go env -w GOPROXY=${GOPROXY} && go env -w GOSUMDB=off && go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o review-manager .

# 使用轻量级的alpine镜像
FROM alpine:latest

# 定义构建参数（第二阶段）
ARG ALPINE_MIRROR=mirrors.aliyun.com

# 安装必要的运行时依赖
RUN sed -i "s#dl-cdn.alpinelinux.org#${ALPINE_MIRROR}#g" /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

# 创建必要的目录
RUN mkdir -p /app/logs /app/config /app/i18n/locales

# 复制二进制文件
COPY --from=builder /app/review-manager /app/
# 复制配置文件
COPY --from=builder /app/config/config.yaml /app/config/
# 复制国际化文件
COPY --from=builder /app/i18n/locales/ /app/i18n/locales/

# 设置工作目录
WORKDIR /app

# 设置权限
RUN chown -R appuser:appuser /app
USER appuser

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O- http://localhost:8080/health || exit 1

# 启动应用
CMD ["./review-manager"]