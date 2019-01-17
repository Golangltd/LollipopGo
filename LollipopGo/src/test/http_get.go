package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// var chandata chan map[int]int

// func init() {
// 	chandata = make(chan map[int]int)
// 	go TTimerAddData()
// 	go TTimeGetData()

// }

// // 压入数据测试
// func TTimerAddData() {

// 	vcount := 1
// 	keycount := 10000

// 	for {
// 		select {
// 		case <-time.After(time.Second * 10):
// 			{
// 				data := make(map[int]int)
// 				data[keycount] = vcount
// 				vcount++
// 				keycount++
// 				chandata <- data
// 			}
// 		}
// 	}
// }

// // 获取数据测试
// func TTimeGetData() {
// 	for {

// 		select {
// 		case <-time.After(time.Second * 1):
// 			{
// 			}
// 		case i := <-chandata:
// 			{
// 				for v := range i {
// 					fmt.Println("-----------------i", i)
// 					fmt.Println("-----------------v", v)
// 					fmt.Println("-----------------vi", i[v])
// 				}
// 			}
// 		}
// 	}
// }

// func main() {

// 	strport := "8892" //  GM 系统操作 -- 修改金币等操作
// 	http.HandleFunc("/GolangLtdGM", IndexHandlerGM)
// 	http.ListenAndServe(":"+strport, nil)
// 	return
// }
