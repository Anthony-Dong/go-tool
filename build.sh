#!/usr/bin/env bash

set -e

# Go的环境变量配置，这个适合与国内的开发者，使用国内的GoProxy
OLDGO111MODULE="$GO111MODULE"
OLDGOPROXY="$GOPROXY"
OLDGOFLAGS="$GOFLAGS"
export GOFLAGS=""
export GO111MODULE=on
export GOPROXY=https://goproxy.cn


# 不开启vendor
go build  -o upload -race -work -v -ldflags "-s" -gcflags "-N -l" cmd/upload-file.go

# 恢复原来的
export GO111MODULE="$OLDGO111MODULE"
export GOPROXY="$OLDGOPROXY"
export GOFLAGS="${OLDGOFLAGS}"

# 创建bin目录
if [ ! -d ./bin ]; then
	mkdir bin
fi

# 移动
if [ -e ./upload ]; then
   mv upload ./bin/
fi

echo 'build finished, out path: ./bin/upload'
