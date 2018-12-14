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

// 游戏的基本的ID的数据信息
func ReadCsv_ConfigFile_GameInfoST_Fun() bool {
	// 获取数据，按照文件
	fileName := "GameInfo.csv"
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

		Infotmp := new(Global_Define.StGameListInfo)
		igame, _ := strconv.Atoi(ss[i][0])
		Infotmp.GameId = uint32(igame)
		Infotmp.GameName = ss[i][1]
		Infotmp.Ip = ss[i][2]
		iport, _ := strconv.Atoi(ss[i][3])
		Infotmp.Port = uint32(iport)
		Infotmp.Type = ss[i][4]
		G_GameInfoST[Infotmp.GameName] = Infotmp
	}

	fmt.Println(G_GameInfoST)
	return true
}
