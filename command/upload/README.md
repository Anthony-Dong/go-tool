# 文件上传工具 - Cli

## 目录：

- [1、特点](#1特点)
- [2、快速开始](#2快速开始)
- [3、配合Typora](#4配合Typora)

## 1、特点

- 利用阿里云Oss，上传图片
- `Typora` 配合使用，写一些markdown，很方便，不需要本地保存图片
- 支持多配置路径，适合上传多个文件
- 个人使用一般是将博客的图片全部上传到阿里云上，个人的一些资料也是，会把url保存住

## 2、快速开始

### 1、下载

```shell
go get -u -v github.com/anthony-dong/go-tool
```

### 2、使用帮助

> ​	配置文件来自于 `go-tool --config 参数`，所以变更配置文件需要指定这个

```shlle
➜ go-tool upload  -h
NAME:
   upload - 文件上传工具

USAGE:
   upload [command options] [arguments...]

OPTIONS:
   --type value, -t value    Set the upload config type (default: "default")
   --file value, -f value    Set the local path of upload file
   --decode value, -d value  Set the upload file name decode method ("uuid"|"base64") (default: "uuid")
   --help, -h                show help (default: false)
```

### 3、简单2

```shell
➜ go-tool upload -f ./go.mod
[GO-TOOL] 2021/01/16 16:52:14 api.go:34: [INFO] [upload] command load config:
{
  "config": "/Users/fanhaodong/project/go-tool/bin/.go-tool.json",
  "decode": "uuid",
  "file": "/Users/fanhaodong/project/go-tool/go.mod",
  "log-level": "debug",
  "type": "default"
}
[GO-TOOL] 2021/01/16 16:52:15 upload-file.go:105: [INFO] [upload] end success, url: https://xxxx.oss-accelerate.xxxx.com/image/2021/1-16/13b00b967xxxxxxxxxxxxxxx.mod
```

### 4、配置文件

```shell
{
  "upload": {
    "default": {
      "access_key_id": "xxxx",
      "access_key_secret": "xxxxxx",
      "endpoint": "oss-accelerate.xxxxx.com",
      "url_endpoint": "xxxx.oss-accelerate.xxxx.com",
      "bucket": "tyut",
      "path_prefix": "image"
    },
    "type-1": {
      "access_key_id": "xxxxxxx",
      "access_key_secret": "xxxxxxx",
      "endpoint": "oss-xxxx.xxxx.com",
      "url_endpoint": "xxx.oss-accelerate.xxxx.com",
      "bucket": "xxxx",
      "path_prefix": "xxxx"
    }
  }
}
```

阿里云Oss端配置大概就是这些：

![image-20200914135934215](https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/42cdf58e904e4dbeac06028639db9d40.png)

如果参数不输入 `-t`，则默认走 `default`标签！，所以一般推荐都设置一个default标签

```shell
➜ go-tool upload -f ./go.mod -t software -d base64
[GO-TOOL] 2021/01/16 16:57:30 api.go:34: [INFO] [upload] command load config:
{
  "config": "/Users/fanhaodong/project/upload-file-cli/bin/.go-tool.json",
  "decode": "base64",
  "file": "/Users/fanhaodong/project/upload-file-cli/go.mod",
  "log-level": "debug",
  "type": "software"
}
[GO-TOOL] 2021/01/16 16:57:34 upload-file.go:105: [INFO] [upload] end success, url: https://anthony-wangpan.oss-accelerate.aliyuncs.com/software/2021/1-16/go.mod
```

## 4、配合Typora

只需要设置如下： 记得`go-tool`写成绝对路径，最后验证一下即可

![image-20210114195430734](https://tyut.oss-accelerate.aliyuncs.com/image/2021/1-14/ec536b08aa054336aaec3a898f203c12.png)


