#!/bin/bash

# 定义语言文件目录
LOCALES_DIR="i18n/locales"

# 处理每个语言文件
for file in "$LOCALES_DIR"/*.yaml; do
    echo "Processing $file..."
    
    # 检查文件是否有custom段
    if ! grep -q "^# custom" "$file"; then
        echo "No custom section found in $file, skipping..."
        continue
    fi
    
    # 创建临时文件
    temp_file=$(mktemp)
    
    # 提取custom段之前的内容
    sed -n '/^# custom/!p;//q' "$file" > "$temp_file"
    
    # 添加custom段标题
    echo "# custom" >> "$temp_file"
    
    # 提取并排序custom段内容
    sed -n '/^# custom/,/^$/p' "$file" | \
        grep -v "^# custom" | \
        grep -v "^$" | \
        sort >> "$temp_file"
    
    # 添加空行分隔
    echo "" >> "$temp_file"
    
    # 添加custom段之后的内容
    sed -n '/^# custom/,$p' "$file" | \
        sed '1,/^$/d' >> "$temp_file"
    
    # 替换原文件
    mv "$temp_file" "$file"
    
    echo "Custom section sorted in $file"
done

echo "All i18n files processed successfully"