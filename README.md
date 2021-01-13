# 文件上传工具 - Cli

## 目录：

- [1、特点](#1特点)
- [2、快速开始](#2快速开始)
- [3、如何使用](#3如何使用)
- [4、配合Typora](#4配合Typora)

## 1、特点

- 利用阿里云Oss，上传图片
- `Typora` 配合使用，写一些markdown，很方便，不需要本地保存图片
- 支持多配置路径，适合上传多个文件
- 个人使用一般是将博客的图片全部上传到阿里云上，个人的一些资料也是，会把url保存住

## 2、快速开始

- 1、直接`go get`

```shell
go get -u -v github.com/anthony-dong/go-tool/cmd
```

- 2、或者下载源码，自己编译

```shell
# 下载代码到本地
wget https://github.com/Anthony-Dong/go-tool/archive/master.zip

# 解压项目
unzip master.zip

# 进入项目目录
cd master

# 编译
make
```

- 3、使用时注意

  > `upload`脚本需要和`aliyun-oss-upload-config.json` 配合使用

```shlle
➜  bin ls | grep upload
aliyun-oss-upload-config.json
upload
```

`aliyun-oss-upload-config.json`内容

1、单配置文件，不需要输入`-t`

```json
{
  "access_key_id": "<access_key_id>",
  "access_key_secret": "<access_key_secret>",
  "endpoint": "oss-accelerate.aliyuncs.com", // 下面图片介绍
  "url_endpoint": "tyut.oss-accelerate.aliyuncs.com",// 下面图片介绍
  "bucket": "tyut", // bucket
  "path_prefix": "image" // 存放的路径，不能在跟路径，必须设置一个
}
```

大概就是这些：

![image-20200914135934215](https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/42cdf58e904e4dbeac06028639db9d40.png)

2、多配置文件

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

## 3、如何使用

- 1、将执行文件导出到环境变量里
- 2、配置`aliyun-oss-upload-config.json` 文件
- 3、直接在目录执行 upload命令，参数是上传的文件路径。生成的文件名称是 `前缀/当前年/当前年-月/uuid.文件格式 `

```shell
➜  /data upload ./Main.java
https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/uuid.java
```

- 4、命令行提示：

```shell
➜  bin upload -h
version:1.0.0
Usage: upload ./Main.java -config=aliyun-oss-upload-config.json

  -config string
    	配置文件位置 (default "aliyun-oss-upload-config.json")
  -h	this help
  -t string
    	-t=default (default "default")
```

## 4、配合Typora

- 1、只需要修改配置即可，十分方便

<img src="https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/e9842bc0463d4900977f4bfe3b43799d.png" alt="image-20200914142037902" style="zoom:50%;" />

- 2、如果没有自动上传，手动点击一下即可

<img src="https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/02a89c4813f3433c8543fb4e5e1db657.png" alt="image-20200914142207048" style="zoom:50%;" />
