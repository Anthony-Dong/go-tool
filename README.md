# Typora - 图片上传Cli



## 1、特点

- 利用阿里云Oss，上传图片
- `Typora` 配合使用，写一些markdown，很方便，不需要本地保存图片

## 2、快速开始

- 1、直接`go get`

```shell
go get github.com/anthony-dong/upload-file-cli/cmd
```

- 2、或者下载源码，自己编译

```shell
wget https://github.com/Anthony-Dong/upload-file-cli/archive/master.zip

unzip master.zip

cd master

./build.sh
```

- 3、使用时注意

  > `upload`脚本需要和`aliyun-oss-upload-config.json` 配合使用

```shlle
➜  bin ls | grep upload
aliyun-oss-upload-config.json
upload
```

`aliyun-oss-upload-config.json`内容

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



## 3、如何使用

1、将执行文件导出到环境变量里

2、配置`aliyun-oss-upload-config.json` 文件

3、直接在目录执行 upload命令，参数是上传的文件路径。生成的文件名称是 `前缀/当前年/当前年-月/uuid.文件格式 `

```shell
➜  /data upload ./Main.java
https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/dc9d4b9a81424e66b73b968fab11a69d.java
```

4、命令行提示：

```shell
➜  aliyun-oss-client-master git:(master) ✗ bin/upload -h        
version:1.0.0
Usage: upload ./Main.java -config=aliyun-oss-upload-config.json

  -config string
        配置文件位置 (default "aliyun-oss-upload-config.json")
  -h    this help
```

