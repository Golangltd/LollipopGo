package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadOne(w http.ResponseWriter, r *http.Request) {
	//判断请求方式
	if r.Method == "POST" {
		//设置内存大小
		r.ParseMultipartForm(32 << 20)
		//获取上传的第一个文件
		file, header, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		//创建上传目录
		os.Mkdir("./upload", os.ModePerm)
		//创建上传文件
		cur, err := os.Create("./upload/" + header.Filename)
		defer cur.Close()
		if err != nil {
			log.Fatal(err)
		}
		//把上传文件数据拷贝到我们新建的文件
		io.Copy(cur, file)
	} else {
		//解析模板文件
		t, _ := template.ParseFiles("./uploadOne.html")
		//输出文件数据
		t.Execute(w, nil)
	}
}

func uploadMore(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//设置内存大小
		r.ParseMultipartForm(32 << 20)
		//获取上传的文件组
		files := r.MultipartForm.File["file"]
		len := len(files)
		for i := 0; i < len; i++ {
			//打开上传文件
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			//创建上传目录
			os.Mkdir("./upload", os.ModePerm)
			//创建上传文件
			cur, err := os.Create("./upload/" + files[i].Filename)
			defer cur.Close()
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(cur, file)
		}
	} else {
		//解析模板文件
		t, _ := template.ParseFiles("./uploadMore.html")
		//输出文件数据
		t.Execute(w, nil)
	}
}

// 通过http://127.0.0.1:9090/uploadOne和http://127.0.0.1:9090/upladMore来测试文件上传。

/*
  策划配置好CSV文件后上传到服务器
  1. 文件上传到制定文件夹
  2. Gm命令通知所有服务器重新加载配置文件，这个时候要求暂停玩家数据的更新，相当于热更新操作！
  3. 账号建权系统的建立
  4. 日志系统
  5. 版本回滚系统等
*/
func main() {
	http.HandleFunc("/uploadMore", uploadMore)
	http.HandleFunc("/uploadOne", uploadOne)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
