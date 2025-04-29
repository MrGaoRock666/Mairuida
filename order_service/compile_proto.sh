#!/bin/bash

# 当前目录为 order_service
PROJECT_ROOT="$(dirname "$(pwd)")"
PROTO_DIR="$(pwd)/pb"
IMPORT_DIR_1="$PROJECT_ROOT/proto"

# 标准库路径
PROTOC_INCLUDE=$(protoc --version >/dev/null 2>&1 && echo "$(dirname $(which protoc))/../include")

echo "开始清理旧的pb文件..."
rm -f "$PROTO_DIR"/*.pb.go

echo "开始编译 order.proto ..."

protoc \
-I "$PROTO_DIR" \
-I "$IMPORT_DIR_1" \
-I "$PROTOC_INCLUDE" \
--go_out=paths=source_relative:"$PROTO_DIR" \
--go-grpc_out=paths=source_relative:"$PROTO_DIR" \
"$PROTO_DIR/order.proto"

# 检查是否成功
if [ $? -eq 0 ]; then
  echo "✅ 编译成功！已生成 order.pb.go 和 order_grpc.pb.go"
else
  echo "❌ 编译失败，请检查 order.proto 是否存在语法错误"
fi