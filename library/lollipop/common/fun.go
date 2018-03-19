/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月5日
*/
package Lcommon

import (
	"encoding/base64"
	"glog-master"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"strings"
)

// 获取路径
func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	CheckErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

// 检测错误
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// html模板设置(渲染模板)
func Assign(html_Path string) (data *template.Template) {
	tmpl, err := template.ParseFiles(html_Path)
	CheckErr(err)
	return tmpl
}

// 保存磁盘的数据的图片处理函数
func SaveFiles(StrPath string, StrBase64Data string, StrPicType string, StrPicName string) bool {
	glog.Info("Entry SaveFiles!")
	glog.Info("SaveFiles  path:" + StrPath)
	//	StrBase64Data ='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'
	SaveFilestmp(StrPath, StrBase64Data, StrPicType, StrPicName)
	// 转换下
	StrBase64Data = strings.Replace(StrBase64Data, "\"", "", -1)
	// 解析
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(StrBase64Data))
	// 转换成png格式的图像，需要导入：_“image/png”
	m, _, _ := image.Decode(reader)
	// 输出到磁盘里:包括路径
	// 文件夹操作
	//dir, _ := os.Getwd() // 获取当前的程序路径
	StrPath = "/var/www/html/res/" + StrPath
	err := os.MkdirAll(StrPath, os.ModePerm) //生成多级目录
	if err != nil {
		glog.Info(err.Error())
		return false
	}
	glog.Info("创建文件夹" + StrPath + "/a/b/c成功")
	StrPicTypetmp := StrPicType
	// 保存数据
	StrPicType = StrPath + "/" + StrPicName + "." + StrPicType
	wt, err := os.Create(StrPicType)
	if err != nil {
		glog.Info("Save Image Error!" + err.Error())
		return false
	}
	// defer wt.Close()
	if wt == nil {
		glog.Info("Save Image Error!  wt is nil!!!")
		return false
	}
	// 转换为jpeg格式的图像，这里质量为30（质量取值是1-100）
	strjsp := "jpg"
	if strings.EqualFold(StrPicTypetmp, strjsp) {
		glog.Info("SaveFiles pic of StrPicType is jpg!")
		jpeg.Encode(wt, m, &jpeg.Options{100})
		wt.Close()
		return true
	}
	// png 图片上传
	strpng := "png"
	glog.Info("save StrPicType:" + StrPicTypetmp)
	if strings.EqualFold(StrPicTypetmp, strpng) {
		glog.Info("SaveFiles pic of StrPicType is png!")
		errpng := png.Encode(wt, m)
		if errpng != nil {
			glog.Info("errpng!", errpng.Error())
		}
		wt.Close()
		return true
	}
	glog.Info("SaveFiles pic of StrPicType is wrong!")
	return false
}

// 保存磁盘的数据的图片处理函数预计装载
func SaveFilestmp(StrPath string, StrBase64Data string, StrPicType string, StrPicName string) bool {
	// 转换下
	StrBase64Data = strings.Replace(StrBase64Data, "\"", "", -1)
	// 解析
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(StrBase64Data))
	// 转换成png格式的图像，需要导入：_“image/png”
	m, _, _ := image.Decode(reader)
	StrPath = "/var/www/html/restmp/" + StrPath
	err := os.MkdirAll(StrPath, os.ModePerm) //生成多级目录
	if err != nil {
		glog.Info(err.Error())
		return false
	}
	StrPicTypetmp := StrPicType
	// 保存数据
	StrPicType = StrPath + "/" + StrPicName + "." + StrPicType
	wt, err := os.Create(StrPicType)
	if err != nil {
		glog.Info("Save Image Error!" + err.Error())
		return false
	}
	// defer wt.Close()
	if wt == nil {
		glog.Info("Save Image Error!  wt is nil!!!")
		return false
	}
	// 转换为jpeg格式的图像，这里质量为30（质量取值是1-100）
	strjsp := "jsp"
	if strings.EqualFold(StrPicTypetmp, strjsp) {
		glog.Info("SaveFiles pic of StrPicType is jpg!")
		//jpeg.Encode(wt, m, &jpeg.Options{100})
		//wt.Close()
		return false
	}
	// png 图片上传
	strpng := "png"
	if strings.EqualFold(StrPicTypetmp, strpng) {
		errpng := png.Encode(wt, m)
		if errpng != nil {
			glog.Info("errpng!", errpng.Error())
		}
		wt.Close()
		return true
	}
	return false
}
