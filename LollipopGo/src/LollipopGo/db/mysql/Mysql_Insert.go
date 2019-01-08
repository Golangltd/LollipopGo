package Mysyl_DB

import (
	_ "LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/player"
	"LollipopGo/LollipopGo/util"
	"Proto/Proto2"
	"database/sql"
	"fmt"
)

/*
   插入数据库数据操作
*/

func insertToDB(db *sql.DB) {
	fmt.Println("insertToDB")
	uid := GetNowtimeMD5()
	nowTimeStr := GetTime()
	stmt, err := db.Prepare("insert t_userinfo set username=?,departname=?,created=?,password=?,uid=?")
	CheckErr(err)
	res, err := stmt.Exec("wangbiao", "研发中心", nowTimeStr, "123456", uid)
	CheckErr(err)
	id, err := res.LastInsertId()
	CheckErr(err)
	if err != nil {
		fmt.Println("插入数据失败")
	} else {
		fmt.Println("插入数据成功：", id)
	}
}

//------------------------------------------------------------------------------
// 玩家游戏的数据
func (this *mysql_db) InsertPlayerGameInfoST2DB(data *Proto2.DB_GameOver) bool {
	// 1.先获取经验，然后累加
	// 2.查表判断
	expdata := this.ReadUserGameExpInfoData(data.OpenID)
	_ = expdata

	//先查询数据,是否存在
	if this.ReadUserGameInfoData(data.OpenID) {
		this.Modefy_PlayerUserGameInfoDataGM(data)
		return false
	}
	//--------------------------------------------------------------------------
	tmptime := util.GetNowUnix_LollipopGo()
	stmt, err := this.STdb.Prepare("insert t_usergameinfo set openid=?,gameid=?,gamelev=?,gameexp=?,gameitem=?,gamescore=?,creattime=?")
	CheckErr(err)
	res, err := stmt.Exec(data.OpenID, data.GameID, data.GameLev, data.GameExp, data.GameScore, data.GameItem, data.GameScore, tmptime)
	CheckErr(err)
	id, err := res.LastInsertId()
	CheckErr(err)
	if err != nil {
		fmt.Println("插入数据失败")
		return false
	} else {
		fmt.Println("插入数据成功：", id)
	}
	//--------------------------------------------------------------------------
	return true
}

//------------------------------------------------------------------------------
// 玩家数据保存
func (this *mysql_db) InsertPlayerST2DB(data *player.PlayerSt) (bool, player.PlayerSt) {
	uid := data.UID
	// 判断是否存在
	bret, bdata := this.ReadUserInfoData(util.Int2str_LollipopGo(uid))
	fmt.Println("数据存在bret！", bret)
	if bret {
		fmt.Println("数据存在！", bdata)
		return false, bdata
	}
	// 获取时间戳
	tmptime := util.GetNowUnix_LollipopGo()
	stmt, err := this.STdb.Prepare("insert t_userinfo set uid=?,openid=?,vip=?,name=?,headurl=?,school=?,sex=?,hallexp=?,coinnum=?,masonrynum=?,mcard=?,constellation=?,medallist=?,createtime=?")
	CheckErr(err)
	res, err := stmt.Exec(data.UID, data.OpenID, data.VIP_Lev, data.Name, data.HeadURL, data.PlayerSchool, data.Sex, data.HallExp, data.CoinNum, data.MasonryNum, data.MCard, data.Constellation, data.MedalList, tmptime)
	CheckErr(err)
	id, err := res.LastInsertId()
	CheckErr(err)
	if err != nil {
		fmt.Println("插入数据失败")
		return false, bdata
	} else {
		fmt.Println("插入数据成功：", id)
	}

	return true, bdata //int(id)
}

//------------------------------------------------------------------------------
