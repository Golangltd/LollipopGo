package main

import (
	"fmt"
	"net/http"
	"time"
)

// 匹配的结构
type Match_player struct {
	UID int
	Lev int
}

// 匹配的chan
var Match_Chan chan *Match_player
var Imax int = 0

// 初始化
func init() {
	Match_Chan = make(chan *Match_player, 100)
	return
}

// 主函数
func main() {

	// 第一个数据：
	idata := &Match_player{
		UID: 1,
		Lev: 6,
	}
	Putdata(idata)

	// 第二个数据：
	idata1 := &Match_player{
		UID: 2,
		Lev: 20,
	}
	Putdata(idata1)

	// 第三个数据：
	idata2 := &Match_player{
		UID: 3,
		Lev: 90,
	}
	Putdata(idata2)

	// 第四个数据：
	idata3 := &Match_player{
		UID: 3,
		Lev: 900,
	}
	Putdata(idata3)
	Putdata(idata3)
	Putdata(idata3)

	// defer close(Match_Chan)
	Imax = len(Match_Chan)
	// 取数据
	//DoingMatch()
	go Sort_timer()

	strport := "8892" //  GM 系统操作 -- 修改金币等操作
	http.HandleFunc("/GolangLtdGM", IndexHandlerGM)
	http.ListenAndServe(":"+strport, nil)

	return
}

// 压入
func Putdata(data *Match_player) {
	// fmt.Print("put:", data, "\t")
	Match_Chan <- data
	// fmt.Print("len:", len(Match_Chan), "\t")
	return
}

// 获取
func DoingMatch() {
	// 全部数据都拿出来
	// data := make(chan map[string]*Match_player, 100)
	// data <- Match_Chan

	for i := 0; i < Imax; i++ {
		if data, ok := <-Match_Chan; ok {
			fmt.Print(data, "\t")
		} else {
			fmt.Print("woring", "\t")
			break
		}

	}
	return
}

//
func Sort_timer() {

	timer := time.NewTimer(time.Second * 1)
	for {
		select {
		case <-timer.C:
			{
				// 获取channel数据的函数。
				DoingMatch()
			}
		}
	}

}

// 排序算法  每8秒排序一次
func Sort_channel() {

	return
}

func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
	// 需要处理 get请求等
}
