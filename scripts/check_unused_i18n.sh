#!/bin/bash

# 定义语言文件目录和项目根目录
LOCALES_DIR="i18n/locales"
PROJECT_ROOT="."
CHECK_PUBLIC=false

# 解析命令行参数
while getopts "p" opt; do
    case $opt in
        p) CHECK_PUBLIC=true ;;
        *) echo "Usage: $0 [-p]" >&2
           exit 1 ;;
    esac
done

# 获取所有custom段的i18n key
get_custom_keys() {
    local file="$1"
    grep -A 1000 "^# custom" "$file" | \
    grep -v "^# custom" | \
    grep -v "^$" | \
    cut -d: -f1 | \
    sed 's/^[[:space:]]*//;s/[[:space:]]*$//'
}

# 获取所有public段的i18n key
get_public_keys() {
    local file="$1"
    grep -A 1000 "^# public" "$file" | \
    grep -v "^# public" | \
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
    
    # 根据参数决定检查哪个段
    if $CHECK_PUBLIC; then
        section="# public"
        section_name="public"
        get_keys_func=get_public_keys
    else
        section="# custom"
        section_name="custom"
        get_keys_func=get_custom_keys
    fi
    
    # 检查文件是否有对应段
    if ! grep -q "$section" "$file"; then
        echo "No $section_name section found in $file, skipping..."
        continue
    fi
    
    # 获取所有key
    keys=$($get_keys_func "$file")
    
    # 检查每个key是否被使用
    unused_keys=()
    for key in $keys; do
        if ! check_key_usage "$key"; then
            unused_keys+=("$key")
        fi
    done
    
    # 输出结果
    if [ ${#unused_keys[@]} -eq 0 ]; then
        echo "All $section_name i18n keys are in use."
    else
        echo "Unused $section_name i18n keys:"
        printf "  - %s\n" "${unused_keys[@]}"
    fi
done