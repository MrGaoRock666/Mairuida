#!/bin/bash

echo "开始安装 grpcurl（从国内镜像下载）..."

TMP_DIR="/tmp/grpcurl_install"
mkdir -p "$TMP_DIR"
cd "$TMP_DIR"

# 国内地址（阿里云 OSS 存储，或者是预上传资源）
URL="https://cdn.jsdelivr.net/gh/long2ice/grpcurl-binary/grpcurl_1.8.7_linux_x86_64"

echo "从国内地址下载 grpcurl..."
curl -LO "$URL"

# 重命名并移动
mv grpcurl_1.8.7_linux_x86_64 grpcurl
chmod +x grpcurl
mv grpcurl /usr/local/bin/

cd ~
rm -rf "$TMP_DIR"

# 验证安装
if command -v grpcurl >/dev/null 2>&1; then
    echo "✅ grpcurl 安装成功！版本：$(grpcurl -version)"
else
    echo "❌ grpcurl 安装失败，请检查网络或权限。"
fi
