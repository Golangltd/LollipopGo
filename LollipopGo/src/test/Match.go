package main

import (
	"fmt"
	"net/http"
	"time"
)

/*
  双人对战匹配 例子
  	1. 适用游戏竞技场、小游戏匹配对战
	2. 技术点：Go语言goroutine间的通信
	3. Go语言实现匹配机制自带buff的，为什么这么说？因为channel本身就是队列的实现，算法优化浑然天成！
*/

//------------------------------------------------------------------------------
/*
   游戏玩家的结构体：
    1. 简单定义几个成员
    2. 正常游戏中一样，只是结构成员不同
*/

type PlayerST struct {
	UID     int
	Name    string
	Lev     int8
	VIP_Lev int8
}

//------------------------------------------------------------------------------
/*
   匹配的chan的定义
*/

var MatchChan chan *PlayerST
var Imax int = 0

//------------------------------------------------------------------------------

/*
  初始化：
    1. 切记 go语言类型中，引用类型必须初始化后才可以使用，map chan slice等
	2. 使用make初始化chan
	3. 创建带有缓冲的chan,因为无缓冲会阻塞玩家排队不合理;如果有不懂的可以文章下面留言
	4. 模拟玩家进入排队chan
*/

func init() {
	// 初始化chan
	MatchChan = make(chan *PlayerST, 100)
	// 玩家A
	player_a := &PlayerST{
		UID:     1,
		Name:    "玩家A",
		Lev:     1,
		VIP_Lev: 0,
	}
	// 玩家B
	player_b := &PlayerST{
		UID:     10,
		Name:    "玩家B",
		Lev:     55,
		VIP_Lev: 0,
	}
	// 玩家C
	player_c := &PlayerST{
		UID:     99,
		Name:    "玩家C",
		Lev:     2,
		VIP_Lev: 1,
	}
	// 放入chan(正常游戏中：客户端发排队消息给服务器，消息带玩家的信息等；服务器接收后同样存入chan)
	MatchChan <- player_a
	MatchChan <- player_b
	MatchChan <- player_c
	go Sort_timer()
}

//------------------------------------------------------------------------------
/*
  取出chan队列里的玩家的数据:
     1. 由于channel的特殊性质，取数据的时候需要注意，不要一次去不取出来
*/

func DoingMatch() {
	Imax = len(MatchChan)
	icount := Imax
	Data := make(map[int]*PlayerST)
	for i := 0; i < Imax; i++ {
		if icount == 1 {
			fmt.Println(MatchChan, "等待匹配")
			continue
		}

		if data, ok := <-MatchChan; ok {
			fmt.Println(data)
			Data[i+1] = data
		} else {
			fmt.Println("woring")
			break
		}
		if icount >= 1 {
			icount--
		}
	}
	if len(Data) > 0 {
		fmt.Println("-------", Data)
	}
}

// 匹配的定时器
func Sort_timer() {
	for {
		select {
		case <-time.After(time.Second * 1):
			{
				DoingMatch()
			}
		}
	}
}

//------------------------------------------------------------------------------

func main() {
	strport := "8888"
	http.HandleFunc("/GolangLtd", IndexHandlerGM)
	http.ListenAndServe(":"+strport, nil)
	return
}

func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
