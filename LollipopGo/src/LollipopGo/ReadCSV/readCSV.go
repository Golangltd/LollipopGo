package csv

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
  功能说明:
      GM 更新操作,游戏中的配置文件
      更新 游戏中的活动,等操作
*/

// GM 操作更新
func ReadCsv_ConfigFile_UpDate_Fun() bool {
	ReadCsv_ConfigFile_GameInfoST_Fun()
	return true
}

func init() {
	// 获取配置列表，游戏列表的数据
	ReadCsv_ConfigFile_GameInfoST_Fun()
}

// 游戏的基本的ID的数据信息
func ReadCsv_ConfigFile_GameInfoST_Fun() bool {
	// 获取数据，按照文件
	fileName := "gamelist.csv"
	fileName = "./csv/" + fileName
	cntb, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("读取配置文件出错!")
		return false
	}
	// 读取文件数据
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	ss, _ := r2.ReadAll()
	sz := len(ss)
	// 循环取数据
	for i := 1; i < sz; i++ {
		Infotmp := new(GameList)
		// igame, _ := strconv.Atoi(ss[i][0])
		// Infotmp.GameId = uint32(igame)
		Infotmp.GameID = ss[i][0]
		Infotmp.GameName = ss[i][1]
		Infotmp.GameICON = ss[i][2]
		Infotmp.IsShow = ss[i][3]
		Infotmp.ShowStartTime = ss[i][4]
		Infotmp.ShowEndTime = ss[i][5]
		Infotmp.IsNewShow = ss[i][6]
		Infotmp.IsHotGame = ss[i][7]
		// 保存数据
		G_GameList[Infotmp.GameName] = Infotmp
	}

	fmt.Println(G_GameInfoST)
	return true
}
