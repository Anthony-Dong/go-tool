// @date : 2020/1/13 15:57
// @author : <a href='mailto:fanhaodong516@qq.com'>Anthony</a>

/**
启动程序
*/

package main

import (
	"fmt"
	"github.com/anthony-dong/aliyun-oss-cli/config"
	"github.com/anthony-dong/aliyun-oss-cli/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	UploadError string = "上传失败"
)

// 这个是B站小电视的图标

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		file, fileInfo, e := r.FormFile("image")
		if e != nil {
			w.WriteHeader(500)
			w.Write([]byte(UploadError))
			return
		}
		image := service.UploadImage(&file, fileInfo.Filename)
		w.Write([]byte(image))
		log.Println(image)
	}
}

func handlerFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if s := r.Method; s == "POST" {
		file, fileInfo, e := r.FormFile("file")
		if e != nil {
			w.WriteHeader(500)
			w.Write([]byte(UploadError))
			return
		}
		image := service.UploadFile(&file, fileInfo.Filename)
		w.WriteHeader(200)
		w.Write([]byte(image))
		log.Println(image)
	}
}

func handlerAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if s := r.Method; s == "POST" {
		file, fileInfo, e := r.FormFile("avatar")
		if e != nil {
			w.WriteHeader(500)
			w.Write([]byte(UploadError))
			return
		}
		image := service.UploadFile(&file, fileInfo.Filename)
		w.WriteHeader(200)
		w.Write([]byte(image))
		log.Println(image)
	}
}

// 服务器地址 : 比如是 47.94.234.232 , 你可以改成 app.4794234232.js .
//text/javascript" src="https://tyut.oss-cn-beijing.aliyuncs.com/static/upload/js/app.4794234232.js

// 本地地址 . 我提供的js文件
// https://tyut.oss-cn-beijing.aliyuncs.com/static/upload/js/app.62f7da10777a76e762dc.js

func main() {
	// 初始化 log
	initLog()

	// 初始化配置
	config.Cinit()

	// 初始化静态资源
	service.InitStatic()

	// 初始化服务器资源
	service.UploadInit()

	// 初始化服务器
	initHttpServer()
}

func initHttpServer() {
	http.HandleFunc("/upload/image", handler)
	http.HandleFunc("/upload/file", handlerFile)
	http.HandleFunc("/upload/avatar", handlerAvatar)
	http.HandleFunc("/img.ico", service.HandIcon)
	http.HandleFunc(service.Static_ICON, service.HandIcon)
	http.HandleFunc(service.STATIC_JS_URL, service.HandlerJS)
	http.HandleFunc(service.STATIC_CSS_URL, service.HandlerCSS)
	http.HandleFunc("/", service.HandIndex)
	fmt.Println("服务启动成功 : ", "http://"+strings.TrimSpace(config.OSSConfig.ServerHost)+":"+strings.TrimSpace(config.OSSConfig.ServerPort))
	e1 := http.ListenAndServe(service.ListenerPort, nil)
	if e1 != nil {
		log.Fatal("服务启动失败 , ", e1)
	}
}

func initLog() {
	format := time.Now().Format("2006-1-2")
	fileName := "upload-file-" + format + ".log"
	// 按照所需读写权限创建文件 , 记住别关闭文件流
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	path, e := filepath.Abs(fileName)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println("日志文件创建成功 : ", path)
	// 完成后延迟关闭
	//设置日志输出到 f
	log.SetOutput(f)
}
