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

```shlle
➜ go-tool upload -h
NAME:
   upload - 文件上传工具

USAGE:
   upload [command options] [arguments...]

OPTIONS:
   --oss_config_file value, -c value   the aliyun oss config file (default: "aliyun-oss-upload-config.json")
   --oss_config_type value, -t value   the aliyun oss config type, default is default (default: "default")
   --file value, -f value              the upload file local path
   --file_name_decode value, -d value  the upload file name decode (default: "uuid")
   --log value                         the log level of print logger (default: "debug")
   --help, -h                          show help (default: false)
```

### 3、简单实用

```shell
➜  f2 go-tool upload -f ./a.txt
[GO-TOOL] 2021/01/14 19:48:51 api.go:35: [INFO] [upload] command load config:
{
  "config_file": "/Users/fanhaodong/go/bin/upload-config.json",
  "config_type": "default",
  "file": "/data/test/f2/a.txt",
  "file_name_decode": "uuid"
}
[GO-TOOL] 2021/01/14 19:48:51 upload-file.go:112: [INFO] [upload] end success, url: https://tyut.oss-accelerate.aliyuncs.com/image/2021/1-14/bd04fa5467fa4f8f88d93fe20558e537.txt
```

### 4、简单配置文件

```shell
{
    "access_key_id": "LTAxxxkPV7oBxxxxxxxx",
    "access_key_secret": "ihxxx2Hkixxxxxxx8cBQNKP5N",
    "endpoint": "oss-xxxx.aliyuncs.com",
    "url_endpoint": "xxxx.oss-accelerate.aliyuncs.com",
    "bucket": "xxxx",
    "path_prefix": "image"
}
```

大概就是这些：

![image-20200914135934215](https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/42cdf58e904e4dbeac06028639db9d40.png)

### 5、多配置文件

```json
{
  "default": {
    "access_key_id": "xxx",
    "access_key_secret": "xxxx",
    "endpoint": "oss-accelerate.aliyuncs.com",
    "url_endpoint": "xxx.oss-accelerate.aliyuncs.com",
    "bucket": "xxx",
    "path_prefix": "image"
  },
  "pdf": {
    "access_key_id": "xxxx",
    "access_key_secret": "xxxx",
    "endpoint": "oss-accelerate.aliyuncs.com",
    "url_endpoint": "xxxx-xxxx.oss-accelerate.aliyuncs.com",
    "bucket": "xxxx-xxxx",
    "path_prefix": "pdf"
  }
}
```

如果参数不输入 `-t`，则默认走 `default`标签！，所以一般推荐都设置一个default标签

```shell
# 文件名称编码 base64, 文件上传配置类别 software
➜ go-tool upload -d base64 -t software -f ./a.txt
[GO-TOOL] 2021/01/14 19:58:52 api.go:35: [INFO] [upload] command load config:
{
  "config_file": "/Users/fanhaodong/go/bin/upload-config.json",
  "config_type": "software",
  "file": "/data/test/f2/a.txt",
  "file_name_decode": "base64"
}
[GO-TOOL] 2021/01/14 19:58:54 upload-file.go:112: [INFO] [upload] end success, url: https://anthony-wangpan.oss-accelerate.aliyuncs.com/software/2021/1-14/a.txt
```

## 4、配合Typora

只需要设置如下： 记得`go-tool`写成绝对路径，最后验证一下即可

![image-20210114195430734](https://tyut.oss-accelerate.aliyuncs.com/image/2021/1-14/ec536b08aa054336aaec3a898f203c12.png)


