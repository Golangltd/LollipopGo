package Mysyl_DB

import (
	"LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/player"
	"LollipopGo/LollipopGo/util"
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
// 玩家数据保存
func (this *mysql_db) InsertPlayerST2DB(data *player.PlayerSt) bool {
	fmt.Println("InsertPlayerST2DBInsertPlayerST2DB")
	uid := data.UID
	// 判断是否存在
	if this.ReadUserInfoData(util.Int2str_LollipopGo(uid)) {
		log.Debug("数据存在！")
		return false
	}
	nowTimeStr := GetTime()
	stmt, err := this.STdb.Prepare("insert t_userinfo set username=?,departname=?,created=?,password=?,uid=?")
	CheckErr(err)
	res, err := stmt.Exec("test", "研发中心", nowTimeStr, "123456", uid)
	CheckErr(err)
	id, err := res.LastInsertId()
	CheckErr(err)
	if err != nil {
		fmt.Println("插入数据失败")
		return false
	} else {
		fmt.Println("插入数据成功：", id)
	}
	return true
}

//------------------------------------------------------------------------------
