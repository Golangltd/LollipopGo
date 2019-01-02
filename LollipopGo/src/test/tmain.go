package main

import (
	"LollipopGo/LollipopGo/util"
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

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	re := InitDSQ(data)
	fmt.Println(re)

	Match_Chan = make(chan *Match_player, 100)
	return
}

func InitDSQ(data1 []int) [4][4]int {

	data := data1
	erdata := [4][4]int{}
	j, k := 0, 0

	// 循环获取
	for i := 0; i < 8*2; i++ {
		// 删除第i个元素
		icount := util.RandInterval_LollipopGo(0, int32(len(data))-1)
		fmt.Println("随机数：", icount)
		//datatmp := data[icount]

		if len(data) == 1 {
			erdata[3][3] = data[0]

		} else {
			//------------------------------------------------------------------
			if int(icount) < len(data) {
				erdata[j][k] = data[icount]
				k++
				if k%4 == 0 {
					j++
					k = 0
				}

				data = append(data[:icount], data[icount+1:]...)
			} else {
				erdata[j][k] = data[icount]
				k++
				if k%4 == 0 {
					j++
					k = 0
				}
				data = data[:icount-1]
			}
			//------------------------------------------------------------------
		}
		fmt.Println("生成的数据", erdata)
	}

	return erdata
}

// 主函数
func main1() {

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
	// DoingMatch()
	//go Sort_timer()

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

	Data := make(map[int]*Match_player)
	// 全部数据都拿出来
	// data := make(chan map[string]*Match_player, 100)
	// data <- Match_Chan
	for i := 0; i < Imax; i++ {
		if data, ok := <-Match_Chan; ok {
			fmt.Print(data, "\t")
			Data[i+1] = data
		} else {
			fmt.Print("woring", "\t")
			break
		}
	}
	// 打印数据保存
	fmt.Println(Data)
	return
}

func Sort_timer() {
	// 控制排队的速度
	timer := time.NewTimer(time.Millisecond * 400)
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

// 数据处理
func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
	// 需要处理 get请求等
}

//------------------------------------------------------------------------------
