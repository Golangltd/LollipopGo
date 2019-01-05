package main

import (
	"cache2go"
	"fmt"
)

var cache *cache2go.CacheTable
var SaveChessData map[int]*GolangLtd

type GolangLtd struct {
	RoomUID   int
	PlayerA   string
	PlayerB   string
	Default   [4][4]int
	ChessData [4][4]int
}

func init() {
	cache = cache2go.Cache("myCache")
	SaveChessData = make(map[int]*GolangLtd)
	return
}

//------------------------------------------------------------------------------

func main() {

	data11 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	data1 := [4][4]int{{17, 17, 17, 17}, {17, 17, 17, 17}, {17, 17, 17, 17}, {17, 17, 17, 17}}
	re := InitDSQ(data11)

	data := &GolangLtd{
		RoomUID:   1,
		PlayerB:   "987654321",
		Default:   data1,
		ChessData: re,
	}

	SaveChessData[1] = data
	// 保存数据
	cache.Add(111, 0, data)
	// 获取数据
	res, err1 := cache.Value(111)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	//--------------------------------------------------------------------------
	fmt.Println("result:", res.Data().(*GolangLtd).RoomUID)
	res.Data().(*GolangLtd).RoomUID = 2
	fmt.Println("result:", res.Data().(*GolangLtd).RoomUID)
	//--------------------------------------------------------------------------
	fmt.Println("result:", res.Data().(*GolangLtd).Default)
	fmt.Println("result:", res.Data().(*GolangLtd).Default[1][2])
	res.Data().(*GolangLtd).Default[1][2] = 18
	fmt.Println("result:", res.Data().(*GolangLtd).Default)
	fmt.Println("result:", res.Data().(*GolangLtd).PlayerA)
	//--------------------------------------------------------------------------
}
