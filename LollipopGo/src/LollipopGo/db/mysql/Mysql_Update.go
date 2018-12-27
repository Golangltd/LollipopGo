package Mysyl_DB

import (
	_ "Proto"
	"Proto/Proto2"
	"database/sql"
	"fmt"
)

/*
   更新数据库数据操作
1 Gm 数据操作 修改数据
*/

// GM 或者数据更新玩家数据操作
func (this *mysql_db) Modefy_PlayerDataGM(uid, itype, number int) bool {
	// 默认修改VIP
	strSql := ""
	if itype == Proto2.MODIFY_COIN {
		strSql = "update t_userinfo_copy set coinnum=? where uid=?"
	} else if itype == Proto2.MODIFY_MASONRY {
		strSql = "update t_userinfo_copy set masonrynum=? where uid=?"
	} else if itype == Proto2.MODIFY_MCARD {
		strSql = "update t_userinfo_copy set mcard=? where uid=?"
	} else if itype == Proto2.MODIFY_LEV {
		strSql = "update t_userinfo_copy set lev=? where uid=?"
		if number > 100 {
			number = 99
		}
	} else if itype == Proto2.MODIFY_VIP_LEV {
		strSql = "update t_userinfo_copy set vip=? where uid=?"
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
