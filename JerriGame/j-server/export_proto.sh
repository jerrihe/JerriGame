#!/bin/bash

# # 检查是否提供目标目录
# if [ $# -lt 1 ]; then
#     echo "Usage: $0 <directory>"
#     exit 1
# fi

# 目标目录
TARGET_DIR="./protocol"

# 检查目标目录是否存在
if [ ! -d "$TARGET_DIR" ]; then
    echo "Error: Directory '$TARGET_DIR' does not exist."
    exit 1
fi

cd $TARGET_DIR

bash export.sh

echo "Shell script execution completed."