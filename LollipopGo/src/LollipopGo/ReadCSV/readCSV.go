package csv

import (
	"LollipopGo/LollipopGo/conf"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	_ "strconv"
	"strings"
)

// 游戏列表
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
		Infotmp := new(conf.GameList)
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
		conf.G_GameList[Infotmp.GameName] = Infotmp
	}

	fmt.Println(conf.G_GameList)
	return true
}

//------------------------------------------------------------------------------

// 游戏列表
func ReadCsv_ConfigFile_BannerInfoST_Fun() bool {
	// 获取数据，按照文件
	fileName := "banner.csv"
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
		Infotmp := new(conf.Banner)
		Infotmp.ADID = ss[i][0]
		Infotmp.PicURL = ss[i][1]
		Infotmp.IsTop = ss[i][2]
		Infotmp.SkipURL = ss[i][3]
		Infotmp.ReMark = ss[i][4]
		// 保存数据
		conf.G_BannerList[Infotmp.ADID] = Infotmp
	}

	fmt.Println(conf.G_BannerList)
	return true
}

//------------------------------------------------------------------------------
