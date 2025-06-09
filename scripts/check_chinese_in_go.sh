#!/bin/bash

# 检测Go文件中是否包含中文
# 用法：./scripts/check_chinese_in_go.sh

set -euo pipefail

# 中文字符的正则表达式（兼容更多grep版本）
CHINESE_CHAR_REGEX='[一-龥]'

# 查找所有.go文件
files=$(find . -name "*.go" -type f)

has_chinese=false

for file in $files; do
    # 检查文件中是否包含中文
    if grep -Pq "$CHINESE_CHAR_REGEX" "$file"; then
        echo "Found Chinese characters in: $file"
        grep -Pn "$CHINESE_CHAR_REGEX" "$file" | while read -r line; do
            echo "  $line"
        done
        has_chinese=true
    fi
done

if [ "$has_chinese" = true ]; then
    echo "Error: Chinese characters found in Go files. Please use i18n instead."
    exit 1
else
    echo "No Chinese characters found in Go files."
    exit 0
fi