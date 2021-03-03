# #######################################################
# Function :Makefile for go                             #
# Platform :All Linux Based Platform                    #
# Version  :1.0                                         #
# Date     :2020-12-17                                  #
# Author   :fanhaodong516@gmail.com                     #
# Usage    :make		   		                        #
# #######################################################

# 项目路径
PROJECT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
# 二进制文件输出位置
GOBUILD_OUT_FILE := bin/go-tool
# 主函数入口
GOBUILD_MAIN_FILE := main.go
# go build参数 -gcflags "-N -l" 参数 -N 来禁用优化，使用 -l 来禁用内联
GOBUILD_ARGS := go build -race -v -ldflags "-s -w"
# go test 相关
GO_TEST_FUNC_NAME := $(test_func)
GO_TEST_PKG_NAME := $(test_pkg)
ifndef GO_TEST_PKG_NAME
	GO_TEST_PKG_NAME := $(shell go list ./...)
endif
# Go的全局的环境变量。GOFLAGS必须清空，防止其他参数干扰
GOFLAGS :=
GO111MODULE := on
GOPROXY := https://goproxy.cn,direct
GOPRIVATE :=
export GO111MODULE
export GOPROXY
export GOPRIVATE
export GOFLAGS

GO_FILES := $(shell find . -name '*.go' | grep -v vendor/ | grep -v _test.go)

# 防止本地文件有重名的问题
.PHONY : all build fmt gofmt goimports golint clean get test testall clear

# make默认启动
all: build

# go build
build: clean fmt
	$(GOBUILD_ARGS) $(GOMOD_VENDOR) -o $(GOBUILD_OUT_FILE) $(GOBUILD_MAIN_FILE)

fmt: gofmt goimports golint clear

gofmt:
	@$(foreach var,$(GO_FILES),echo gofmt -d -w  $(var);gofmt -d -w  $(var);)

goimports:
	@if [ ! -d $(PROJECT_DIR)/bin ]; then mkdir -p $(PROJECT_DIR)/bin; fi
	@if [ ! -e $(PROJECT_DIR)/bin/goimports ]; then curl -o $(PROJECT_DIR)/bin/goimports https://anthony-wangpan.oss-accelerate.aliyuncs.com/software/2020/12-29/788bd0e30957478488d4159859d29a0e && chmod 0744 $(PROJECT_DIR)/bin/goimports; fi
	@$(foreach var,$(GO_FILES),echo goimports -d -w $(var);$(PROJECT_DIR)/bin/goimports -d -w $(var);)

govet:
	@$(foreach var,$(GO_FILES),echo go vet $(GOMOD_VENDOR) $(var);go vet $(GOMOD_VENDOR) $(var);)

golint:
	@if [ ! -d $(PROJECT_DIR)/bin ]; then mkdir -p $(PROJECT_DIR)/bin; fi
	@if [ ! -e $(PROJECT_DIR)/bin/golint ]; then curl -o $(PROJECT_DIR)/bin/golint https://anthony-wangpan.oss-accelerate.aliyuncs.com/software/2020/12-30/6fda119141b84c77b0924e9d140704d0 && chmod 0744 $(PROJECT_DIR)/bin/golint; fi
	@$(foreach var,$(GO_FILES),echo golint $(var);$(PROJECT_DIR)/bin/golint $(var);)

clean:
	$(RM) -r $(GOBUILD_OUT_FILE) coverage.txt

get:
	go get -u -v $(import)
	go mod download

test: clean
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic -run $(GO_TEST_FUNC_NAME) $(GO_TEST_PKG_NAME)
	go tool cover -html=coverage.txt

testall: clean
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt

clear:
	go run clear.go