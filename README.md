# Go-Tool

> ​	主要是使用GO写的一些工具，集成了一些通用的工具包

## 快速开始

### 1、下载

```shell
go get -u -v github.com/anthony-dong/go-tool
```

### 2、使用

可以看到目前只有 `upload`命令，支持全局配置`--config` 和 `--log-level`  ，任何命令都通用，使用方式`go-tool --config /data/config.json upload -f ./a.txt`

```shell
➜ bin/go-tool 
NAME:
   go-tool - A go toll cli application for go

USAGE:
   go-tool [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

COMMANDS:
   upload   文件上传工具
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value     Location of client config files (default: "/Users/fanhaodong/project/upload-file-cli/bin/.go-tool.json")
   --log-level value  Set the logging level (warn|info|debug|fatal|error) (default: "debug")
   --help, -h         show help (default: false)
   --version, -v      print the version (default: false)
```

## 如何贡献代码
1、所有的命令都需要实现这个接口，接口需要实现在以下位置`command/{cli-command}`，可以参考：[upload的实现](./command/upload/cli.go)
```go
type Command interface {
	Run(context *cli.Context) error
	Flag() []cli.Flag
	InitConfig(context *cli.Context, config api.CommonConfig) ([]byte, error)
}
```
2、注册命令,可以在[main.go](./main.go) 中实现
```go
Commands: []*cli.Command{
   command.NewCommand({cli-command}, {cli-desc}, {cli-pck}.NewCommand()),
},
```


## [Go-Upload](./command/upload)

>   使用Aliyun Oss上传文件，目前方便使用，且集成了Typora使用

```shell
➜  upload-file-cli git:(master) ✗ bin/go-tool upload -f ./go.mod
[GO-TOOL] 2021/01/16 16:52:14 api.go:34: [INFO] [upload] command load config:
{
  "config": "/Users/fanhaodong/project/upload-file-cli/bin/.go-tool.json",
  "decode": "uuid",
  "file": "/Users/fanhaodong/project/upload-file-cli/go.mod",
  "log-level": "debug",
  "type": "default"
}
[GO-TOOL] 2021/01/16 16:52:15 upload-file.go:105: [INFO] [upload] end success, url: https://tyut.oss-accelerate.aliyuncs.com/image/2021/1-16/13b00b9672aa43e681a1b5df3bfaf2c8.mod
```
## Go-Wrk

>  压测Http接口工具

```shell
➜  go-tool git:(master)  bin/go-tool wrk -h
NAME:
   wrk - 压测接口工具

USAGE:
   wrk [command options] [arguments...]

OPTIONS:
   --duration value, -d value     Set the request duration for wrk ("1ms"|"1s"|"1m"|"1h") (default: 0s)
   --connections value, -c value  Connections to keep open, must greater than threads (default: 12)
   --threads value, -t value      Number of threads to use (default: 12)
   --method value, -m value       Set the http request method method (default: "GET")
   --url value, -u value          Set the request url
   --body value, -b value         Set the request body
   --header value, -H value       Set the request header
   --timeout value                Socket request timeout (default: 5s)
   --help, -h                     show help (default: false)
```

例如请求某个接口：

```shell
➜  go-tool git:(master) ✗  bin/go-tool wrk -t 10 -c 10 -d 10s -b hello -H content-type:application/json -u http://localhost:9999/test
[GO-TOOL] 2021/05/20 14:49:40 api.go:34: [INFO] [wrk] command load config:
{
  "body": "hello",
  "connections": 10,
  "duration": "10s",
  "header": {
    "Content-Type": [
      "application/json"
    ]
  },
  "method": "GET",
  "threads": 10,
  "timeout": "5s",
  "url": "http://localhost:9999/test"
}
==========请求时长分布图==============
1ms  257038  99.74%
5ms  257681  99.99%
10ms  257704  100.00%
50ms  257709  100.00%
100ms  257709  100.00%
500ms  257709  100.00%
1000ms  257709  100.00%
5000ms  257709  100.00%
==========请求次数分布图==============
请求总次数: 257709
请求次数/s: 25770
==========请求体吞吐量==============
请求吞吐量(kb): 0kb
请求每秒吞吐量(kb/s): 0kb
程序结束，一共花费 10.001s
```

## Go-Orm

> ​	未加入，期待

## [Go-Markdown](./command/markdown)

> ​	自动生成博客的README文件

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

