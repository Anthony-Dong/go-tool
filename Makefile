# #######################################################
# Function :Makefile for go                             #
# Platform :All Linux Based Platform                    #
# Version  :1.0                                         #
# Date     :2020-12-17                                  #
# Author   :fanhaodong516@gmail.com                     #
# Usage    :make		   		                        #
# #######################################################

# dir
PROJECT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# go test 相关
GO_TEST_PKG_NAME := $(shell go list ./...)

# Go env
export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct
export GOPRIVATE :=
export GOFLAGS :=

# 防止本地文件有重名的问题
.PHONY : all init build fmt clean  clean test testall clear help

# make默认启动
all: build

init: ## init
	go mod download
	@if [ ! -e $(shell go env GOPATH)/bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; fi

build: clean fmt ## build
	 go build -v -ldflags "-s -w"  -o bin/go-tool main.go

fmt:
	$(shell go env GOPATH)/bin/golangci-lint run --fix --skip-files _test.go$$  --disable-all --enable govet \
 	--enable gofmt \
 	--enable goimports \
 	--enable godot \
 	--enable whitespace \
 	--enable gci
clean:
	$(RM) -r bin/go-tool coverage.txt

get:
	go get -u -v $(import)
	go mod download

test: clean ## test
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic -run $(test_func) $(test_pkg)
	go tool cover -html=coverage.txt

testall: clean ## test all
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt

clear: ## 清空敏感文字
	go build -v -o clear-tool clear/clear.go
	./clear-tool

help: ## 帮助
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)