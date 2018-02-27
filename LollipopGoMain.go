package main

import (
	"fmt"
	"html/template"
	"os"
)

// 初始化函数
func init() {
	fmt.Println("Entry Init()")
	return
}

// 主函数
func main() {
	fmt.Println("Entry Main()")
	type person struct {
		Id      int
		Name    string
		Country string
	}
	//
	liumiaocn := person{Id: 1001, Name: "liumiaocn", Country: "China"}
	fmt.Println("liumiaocn = ", liumiaocn)
	// 输出到页面？？ 直接调用
	tmpl, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Println("Error happened..")
	}
	tmpl.Execute(os.Stdout, liumiaocn)
	return
}
