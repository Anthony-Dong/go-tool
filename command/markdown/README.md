# 自动生成博客README - Cli

## 目录：

- [1、特点](#1特点)
- [2、快速开始](#2快速开始)
- [3、配合Typora](#4配合Typora)

## 1、特点

- 生成博客的README文件
- 本人博客使用的`Typora`，所以很容易构建

## 2、快速开始

### 1、下载

```shell
go get -u -v github.com/anthony-dong/go-tool
```

### 2、使用帮助

> ​	配置文件来自于 `go-tool --config 参数`，所以变更配置文件需要指定这个

```shlle
➜  go-tool markdown -h
NAME:
   markdown - 生成markdown项目的README文件

USAGE:
   markdown [command options] [arguments...]

OPTIONS:
   --dir value, -d value       The markdown project dir
   --template value, -t value  The markdown template file path, go template struct: &{UrlPath: Total:0}
   --help, -h                  show help (default: false)
```

### 3、简单操作

```shell
➜  go-tool --config .config/go-tool.json  markdown --dir /Users/fanhaodong/note/note --template .config/README-template.md
[GO-TOOL] 2021/01/26 16:29:42 api.go:34: [INFO] [markdown] command load config:
{
  "config": "/Users/fanhaodong/note/note/.config/go-tool.json",
  "dir": "/Users/fanhaodong/note/note",
  "log-level": "debug",
  "template": "/Users/fanhaodong/note/note/.config/README-template.md"
}
[GO-TOOL] 2021/01/26 16:29:42 markdown.go:72: [INFO] Get ReadmeFileInfo success, UrlPath: Not show, Total: 452
[GO-TOOL] 2021/01/26 16:29:42 markdown.go:80: [INFO] Create /Users/fanhaodong/note/note/README.md file success !!
[GO-TOOL] 2021/01/26 16:29:42 markdown.go:86: [INFO] New parser success, template file: /Users/fanhaodong/note/note/.config/README-template.md
[GO-TOOL] 2021/01/26 16:29:42 markdown.go:91: [INFO] Write README file success !!
```

### 4、配置文件

> ​	无

### 5、博客的Makefile文件

> ​	博客目录需要有 `.config/go-tool.json` 文件和 `.config/README-template.md` 文件

- `Makefile`文件

```makefile
# #######################################################
# Function :Makefile for go                             #
# Platform :All Linux Based Platform                    #
# Version  :1.0                                         #
# Date     :2021-01-26                                  #
# Author   :fanhaodong516@gmail.com                     #
# Usage    :make		   		                        #
# #######################################################

# 项目路径
PROJECT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# =========构建工具========
# BIN 目录
BIN_DIR := $(PROJECT_DIR)/bin
# GO_TOOL 文件
GO_TOOL := $(BIN_DIR)/go-tool
# GO_TOOL 的GIT地址
GO_TOOL_GIT := github.com/anthony-dong/go-tool
# Go 编译的环境，由于依赖于一些google的包，所以需要使用goproxy.io
GO_ENV := GOPROXY="https://goproxy.io,direct" GOBIN="$(BIN_DIR)"


# ===========构建配置文件===========
GO_TOOL_CONFIG_FILE := $(PROJECT_DIR)/.config/go-tool.json

.PHONY: $(GO_TOOL_CONFIG_FILE) all tool upload

all: tool $(GO_TOOL_CONFIG_FILE)
	$(GO_TOOL) --config $(GO_TOOL_CONFIG_FILE)  markdown --dir $(PROJECT_DIR) --template $(PROJECT_DIR)/.config/README-template.md

tool:
	@if [ ! -d $(BIN_DIR) ]; then mkdir -p $(PROJECT_DIR)/bin; fi
	@if [ ! -e $(GO_TOOL) ]; then $(GO_ENV) go get -u -v $(GO_TOOL_GIT); fi

$(GO_TOOL_CONFIG_FILE):
	-mkdir -p $(shell dirname $(GO_TOOL_CONFIG_FILE))

# ======go-tool upload -f ./README.md -t software -d base64 ==============
UPLOAD :=
ifdef file
	UPLOAD := $(UPLOAD) -f $(file)
endif

ifdef type
	UPLOAD := $(UPLOAD) -t $(type)
endif

ifdef decode
	UPLOAD := $(UPLOAD) -d $(decode)
endif

upload: tool $(GO_TOOL_CONFIG_FILE)
	$(GO_TOOL) --config $(GO_TOOL_CONFIG_FILE) upload $(UPLOAD)

clean:
	$(RM) -r $(BIN_DIR)
```

- `.config/README-template.md` 

  >使用的Go的模版渲染的`README`文件，目前只有`{{.Total}}`和 `{{.UrlPath}}` 元素

```markdown
# 个人笔记

## 1、介绍

- 本人所有的笔记都在里面，使用Git仓库作为平台搭建的，为了查询速度，本人添加了文件目录的选项，方便直接查询，合计{{.Total}}篇
- 有可能涉及到工作公司的一些隐私URL，所以这个不会作为public的仓库
- 目前仓库所有文件都是文本，图片上传依赖于本人写的OSS插件配合Typora ，有兴趣的可以下载OSS插件：[https://github.com/Anthony-Dong/go-tool](https://github.com/Anthony-Dong/go-tool) ，记得star一下
- 本人Github地址：[https://github.com/Anthony-Dong](https://github.com/Anthony-Dong)
- 本人Leetcode项目地址：[https://github.com/Anthony-Dong/leetcode](https://github.com/Anthony-Dong/leetcode)

## 2、文件目录

{{.UrlPath}}
```

