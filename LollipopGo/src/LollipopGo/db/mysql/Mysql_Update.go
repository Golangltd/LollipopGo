package Mysyl_DB

import (
	"LollipopGo/LollipopGo/util"
	"Proto/Proto2"
	"database/sql"
	"fmt"
)

//------------------------------------------------------------------------------
/*
   更新邮件，重置系统邮件
*/

func (this *mysql_db) Modefy_AdminGameEmailInfoDataGM() bool {

	strSql := "update t_adminemail set state=? where id > 0"
	// 修改数据
	stmt, err := this.STdb.Prepare(strSql)
	defer stmt.Close()
	CheckErr(err)
	res, err := stmt.Exec(0)
	affect, err := res.RowsAffected()
	fmt.Println("更新数据：", affect)
	CheckErr(err)

	return true
}

//------------------------------------------------------------------------------
/*
   更新数据，DSQ的修改

*/
func (this *mysql_db) Modefy_PlayerUserGameInfoDataGM(data *Proto2.DB_GameOver, gamelev int) bool {

	strSql := "update t_usergameinfo set gameid=?,gamelev=?,gameexp=?,gameitem=?,gamescore=?,creattime=? where openid=?"
	// 修改数据
	stmt, err := this.STdb.Prepare(strSql)
	CheckErr(err)
	tmptime := util.GetNowUnix_LollipopGo()
	res, err := stmt.Exec(data.GameID, gamelev, data.GameExp, data.GameScore, data.GameItem, data.GameScore, tmptime, data.OpenID)
	affect, err := res.RowsAffected()
	fmt.Println("更新数据：", affect)
	CheckErr(err)

	return true
}

//------------------------------------------------------------------------------
/*
   更新数据库数据操作
1 Gm 数据操作 修改数据
*/

// GM 或者数据更新玩家数据操作
func (this *mysql_db) Modefy_PlayerDataGM(uid, itype, number int) bool {
	// 默认修改VIP
	strSql := ""
	if itype == Proto2.MODIFY_COIN {
		strSql = "update t_userinfo set coinnum=? where uid=?"
	} else if itype == Proto2.MODIFY_MASONRY {
		strSql = "update t_userinfo set masonrynum=? where uid=?"
	} else if itype == Proto2.MODIFY_MCARD {
		strSql = "update t_userinfo set mcard=? where uid=?"
	} else if itype == Proto2.MODIFY_LEV {
		strSql = "update t_userinfo set lev=? where uid=?"
		if number > 100 {
			number = 99
		}
	} else if itype == Proto2.MODIFY_VIP_LEV {
		strSql = "update t_userinfo set vip=? where uid=?"
		if number > 100 {
			number = 99
		}
	}

	if len(strSql) == 0 {
		return false
	}

	// 修改数据
	stmt, err := this.STdb.Prepare(strSql)
	CheckErr(err)
	res, err := stmt.Exec(number, uid)
	affect, err := res.RowsAffected()
	fmt.Println("更新数据：", affect)
	CheckErr(err)

	return true
}

//------------------------------------------------------------------------------
func UpdateDB(db *sql.DB, uid string) {
	stmt, err := db.Prepare("update userinfo set username=? where uid=?")
	CheckErr(err)
	res, err := stmt.Exec("zhangqi", uid)
	affect, err := res.RowsAffected()
	fmt.Println("更新数据：", affect)
	CheckErr(err)
}
