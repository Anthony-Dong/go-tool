# Go-Tool

> ​	主要是使用GO写的一些工具，集成了一些通用的工具包

## 快速开始

```shell
go get -u -v github.com/anthony-dong/go-tool
```


## [Go-Upload](./command/upload)

>   使用Aliyun Oss上传文件，目前方便使用，且集成了Typora使用

```shell
➜  ~ go-tool upload -f /data/test/f1/a 
[GO-TOOL] 2021/01/14 16:33:41 api.go:35: [INFO] [upload] command load config:
{
  "file": "/data/test/f1/a",
  "file_name_decode": "uuid",
  "oss_config_file": "/Users/fanhaodong/go/bin/aliyun-oss-upload-config.json",
  "oss_config_type": "default"
}
[GO-TOOL] 2021/01/14 16:33:46 upload-file.go:112: [INFO] [upload] end success, url: https://tyut.oss-accelerate.aliyuncs.com/image/2021/1-14/0ee1a9e9647b43a38ce3982d5652xxxx
```
## Go-Wrk

>  压测Http接口工具

## Go-Orm