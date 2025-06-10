#!/bin/bash

# 定义语言文件目录
LOCALES_DIR="i18n/locales"
SECTION="# custom"
SECTION_NAME="custom"

# 解析命令行参数
while getopts "p" opt; do
    case $opt in
        p) SECTION="# public"
           SECTION_NAME="public" ;;
        *) echo "Usage: $0 [-p]" >&2
           exit 1 ;;
    esac
done

# 处理每个语言文件
for file in "$LOCALES_DIR"/*.yaml; do
    echo "Processing $file..."
    
    # 检查文件是否有对应段
    if ! grep -q "$SECTION" "$file"; then
        echo "No $SECTION_NAME section found in $file, skipping..."
        continue
    fi
    
    # 创建临时文件
    temp_file=$(mktemp)
    
    # 提取段之前的内容
    sed -n "/$SECTION/!p;//q" "$file" > "$temp_file"
    
    # 添加段标题
    echo "$SECTION" >> "$temp_file"
    
    # 提取并排序段内容
    sed -n "/$SECTION/,/^$/p" "$file" | \
        grep -v "$SECTION" | \
        grep -v "^$" | \
        sort >> "$temp_file"
    
    # 添加空行分隔
    echo "" >> "$temp_file"
    
    # 添加段之后的内容
    sed -n "/$SECTION/,\$p" "$file" | \
        sed '1,/^$/d' >> "$temp_file"
    
    # 替换原文件
    mv "$temp_file" "$file"
    
    echo "$SECTION_NAME section sorted in $file"
done

echo "All i18n files processed successfully"