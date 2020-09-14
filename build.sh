#!/usr/bin/env bash

set -e

# 如果本地没有vendor目录，会选择执行vendor
if [ ! -d ./vendor ]; then
	go mod vendor
fi

# -mod=vendor 强制开启只去vendor目录加载，不会去自动拉去
go build -mod=vendor -o upload -race -work -v -ldflags "-s" -gcflags "-N -l" cmd/upload-file.go

# copy op-report-go
if [ ! -d ./bin ]; then
	mkdir bin
fi

if [ -e ./upload ]; then
   mv upload ./bin/
fi

echo 'build finished'
