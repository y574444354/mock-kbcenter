#!/bin/bash

# 获取镜像名称
# 1. 尝试从项目根目录的 .env 文件读取
if [ -f .env ]; then
  source .env
  if [ -n "$DOCKER_IMAGE" ]; then
    IMAGE_NAME="${DOCKER_IMAGE%:*}"
  elif [ -n "$APP_NAME" ]; then
    IMAGE_NAME="zgsm/${APP_NAME}"
  fi
fi

# 2. 如果未从 .env 获取，使用默认值
IMAGE_NAME=${IMAGE_NAME:-"zgsm/go-webserver"}

# 使用 Makefile 构建 Docker 镜像
echo "Building Docker image using Makefile..."
make docker-build

# 如果有提供 TAG 参数，则打标签并推送
if [ -n "$1" ]; then
  TAG=$1
  echo "Tagging image with $TAG..."
  docker tag $IMAGE_NAME:latest $IMAGE_NAME:$TAG
  echo "Pushing tagged image..."
  docker push $IMAGE_NAME:$TAG
fi

# 推送 latest 镜像
echo "Pushing latest image..."
docker push $IMAGE_NAME:latest

echo "Done."