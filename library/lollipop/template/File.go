/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月2日
*/
package Ltemplate

import (
	"os"
	"os/exec"
	"strings"
)

// 错误提示
var Error_Path string = getCurrentPath() + "error.html"

// 获取路径
func getCurrentPath() string {
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
