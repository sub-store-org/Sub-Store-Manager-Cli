#!/bin/bash

# 检测是否已安装 Docker
if ! command -v docker &> /dev/null; then
    echo "系统未安装 Docker Engine，由于 ssm 依赖于 Docker Engine，请先安装 Docker。"
    exit 1
fi

# 检测是否已安装 curl
if ! command -v curl &> /dev/null; then
    echo "未安装 curl，请先安装 curl。"
    exit 1
fi

# 请将以下变量替换为你要安装的 GitHub Release 信息
GITHUB_REPO_OWNER="DesnLee"
GITHUB_REPO_NAME="Sub-Store-Manager-Cli"
BINARY_NAME="ssm"  # 二进制文件名，不包含后缀
BIN_DIRECTORY="/usr/local/bin"  # 安装到的目标目录
RELEASE_TAG=$(curl -s "https://api.github.com/repos/${GITHUB_REPO_OWNER}/${GITHUB_REPO_NAME}/releases/latest" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)

# 根据操作系统类型和架构确定文件后缀
SYS_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
OS=""
case "$SYS_OS" in
    linux*)  OS="linux" ;;
    darwin*) OS="mac" ;;
    msys*)   OS="windows" ;;
    *)       echo "不支持的操作系统: ${SYS_OS}" && exit 1 ;;
esac

ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

BINARY_FILENAME="${BINARY_NAME}_${OS}_${ARCH}"

echo "准备安装 ${BINARY_NAME} ${RELEASE_TAG}..."
echo "检测系统环境为 ${OS} ${ARCH}..."
echo "下载 ${BINARY_FILENAME}..."

# 下载 GitHub Release 文件
DOWNLOAD_URL="https://github.com/${GITHUB_REPO_OWNER}/${GITHUB_REPO_NAME}/releases/download/${RELEASE_TAG}/${BINARY_FILENAME}"
TMP_DIR=$(mktemp -d)
curl -L -o "${TMP_DIR}/${BINARY_FILENAME}" "${DOWNLOAD_URL}"

# 将下载的文件移动到 bin 目录
chmod +x "${TMP_DIR}/${BINARY_FILENAME}"
sudo mv "${TMP_DIR}/${BINARY_FILENAME}" "${BIN_DIRECTORY}/${BINARY_NAME}"

# 清理临时文件
rm -rf "${TMP_DIR}"

# 验证安装是否成功
if command -v "${BINARY_NAME}" &> /dev/null; then
    echo "${BINARY_NAME} 安装成功！可以使用 ssm -h 查看帮助"
else
    echo "安装失败，请检查问题并重试。"
fi
