package csv

import (
	"LollipopGo/LollipopGo/conf"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"
)

// 读取服务器列表
func ReadCsv_ConfigFile_ServerListInfoST_Fun() bool {
	fileName := "serverlist.csv"
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
	for i := 1; i < sz; i++ {
		Infotmp := new(conf.ServerList)
		Infotmp.ID = ss[i][0]
		Infotmp.Name = ss[i][1]
		Infotmp.IP_Port = ss[i][2]
		conf.G_ServerList[Infotmp.ID] = Infotmp
	}

	return true
}

// 斗兽棋游戏经验列表
func ReadCsv_ConfigFile_DSQGameInfoST_Fun() bool {
	fileName := "Animal_exp.csv"
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
	for i := 1; i < sz; i++ {
		Infotmp := new(conf.DSQ_Exp)
		Infotmp.Level = ss[i][0]
		Infotmp.Exp = ss[i][1]
		conf.DSQGameExp[Infotmp.Level] = Infotmp
	}

	return true
}

// 游戏列表New
func ReadCsv_ConfigFile_GameInfoST_FunNew() bool {
	// 获取数据，按照文件
	fileName := "gamelistnew.csv"
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
		Infotmp := new(conf.GameListNew)
		Infotmp.GameID = ss[i][0]
		Infotmp.Name = ss[i][1]
		Infotmp.IconPath = ss[i][2]
		Infotmp.IsShelves = ss[i][3]
		Infotmp.StartTime = ss[i][4]
		Infotmp.EndTime = ss[i][5]
		Infotmp.IsNewest = ss[i][6]
		Infotmp.IsHot = ss[i][7]
		Infotmp.ResPath = ss[i][8]
		// 保存数据
		conf.G_GameListNew[Infotmp.GameID] = Infotmp
		// 保存数据更新数据
		M_CSV.Set(Infotmp.GameID, 0)
	}
	return true
}

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
		conf.G_GameList[Infotmp.GameID] = Infotmp
		// 保存数据更新数据
		M_CSV.Set(Infotmp.GameID, 0)
	}
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

// 房间列表
// 数据在网关服务器 --- update
func ReadCsv_ConfigFile_RoomListST_Fun() bool {
	// 获取数据，按照文件
	fileName := "roomlist.csv"
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
	roomidtemp := ""
	// 循环取数据
	for i := 1; i < sz; i++ {
		Infotmp := new(conf.RoomList)
		Infotmp.RoomID = ss[i][0]
		Infotmp.NeedPiece = ss[i][1]
		Infotmp.NeedLev = ss[i][2]
		Infotmp.Desc = ss[i][3]
		Infotmp.SysPiece = ss[i][4]
		Infotmp.WinReward = ss[i][5]
		Infotmp.IsTop = ss[i][6]
		Infotmp.TypeICON = ss[i][7]

		s := string([]byte(Infotmp.RoomID)[:5])
		if len(roomidtemp) == 0 {
			roomidtemp = s
			conf.RoomListData[Infotmp.RoomID] = Infotmp
			conf.RoomListDatabak[Infotmp.RoomID] = Infotmp

		} else {
			if roomidtemp == s {
				conf.RoomListData[Infotmp.RoomID] = Infotmp
				conf.RoomListDatabak[Infotmp.RoomID] = Infotmp
				fmt.Println("+++++++++", conf.RoomListData)
				// 仅仅有一个游戏的时候
				if i == sz-1 {
					conf.G_RoomList[roomidtemp] = conf.RoomListData
				}
			} else {
				// 保存数据
				conf.G_RoomList[roomidtemp] = conf.RoomListData
				roomidtemp = s
				conf.RoomListData = make(map[string]*conf.RoomList)

			}
		}
	}
	fmt.Println(conf.G_RoomList["10001"])
	return true
}

//------------------------------------------------------------------------------
