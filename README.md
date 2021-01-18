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

## Go-Orm