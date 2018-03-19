/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月3日
*/
package LollipopGo

import (
	"LollipopGo/library/lollipop/log"
	"fmt"
)

// 配置第三方包的配置文件
// 可以是否打开
func init() {
	fmt.Println("Entry init!!!")
	// 日志开启自动
	if true {
		flag.Set("alsologtostderr", "true") // 日志写入文件的同时，输出到stderr
		flag.Set("log_dir", "./log")        // 日志文件保存目录.执行文件的跟目录
		flag.Set("v", "3")                  // 配置V输出的等级。
		flag.Parse()
	}
	return
}
