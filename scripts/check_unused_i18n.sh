#!/bin/bash

# 定义语言文件目录和项目根目录
LOCALES_DIR="i18n/locales"
PROJECT_ROOT="."

# 获取所有custom段的i18n key
get_custom_keys() {
    local file="$1"
    grep -A 1000 "^# custom" "$file" | \
    grep -v "^# custom" | \
    grep -v "^$" | \
    cut -d: -f1 | \
    sed 's/^[[:space:]]*//;s/[[:space:]]*$//'
}

# 检查key是否在go文件中出现
check_key_usage() {
    local key="$1"
    # 直接使用key进行搜索，不转义点号
    grep -r --include="*.go" -F "$key" "$PROJECT_ROOT" > /dev/null
    return $?
}

# 处理每个语言文件
for file in "$LOCALES_DIR"/*.yaml; do
    echo "Checking $file..."
    
    # 检查文件是否有custom段
    if ! grep -q "^# custom" "$file"; then
        echo "No custom section found in $file, skipping..."
        continue
    fi
    
    # 获取所有custom key
    keys=$(get_custom_keys "$file")
    
    # 检查每个key是否被使用
    unused_keys=()
    for key in $keys; do
        if ! check_key_usage "$key"; then
            unused_keys+=("$key")
        fi
    done
    
    # 输出结果
    if [ ${#unused_keys[@]} -eq 0 ]; then
        echo "All custom i18n keys are in use."
    else
        echo "Unused custom i18n keys:"
        printf "  - %s\n" "${unused_keys[@]}"
    fi
done