# #######################################################
# Function :Makefile for go                             #
# Platform :All Linux Based Platform                    #
# Version  :1.0                                         #
# Date     :2020-12-17                                  #
# Author   :fanhaodong516@gmail.com                     #
# Usage    :make help		   		                    #
# #######################################################

# dir
PROJECT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# go test
GO_TEST_PKG_NAME := $(shell go list ./...)

# go env
export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct
export GOPRIVATE :=
export GOFLAGS :=

# PHONY
.PHONY : all init build fmt check clean test testall clear help

all: build ## Let's go!

init: ## init project and init env
	go mod download
	@if [ ! -e $(shell go env GOPATH)/bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; fi

build: fmt ## build this project
	 go build -v -ldflags "-s -w"  -o bin/go-tool main.go

fmt: ## fmt add auto save
	@$(shell go env GOPATH)/bin/golangci-lint run --fix --skip-files _test.go$$  --disable-all --enable govet \
 	--enable gofmt \
 	--enable goimports \
 	--enable godot \
 	--enable whitespace \
 	--enable gci \

check: ## check this project bugs
	@$(shell go env GOPATH)/bin/golangci-lint run --fix --skip-files _test.go$$ --disable-all \
	--enable errcheck

clean: ## clear not useful file
	$(RM) -r bin coverage.txt clear-tool

test: clean ## go test
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic -run $(test_func) $(test_pkg)
	go tool cover -html=coverage.txt

testall: clean ## go test for the package
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt

clear: ## the clear tool
	go build -v -o clear-tool clear/clear.go
	./clear-tool

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)