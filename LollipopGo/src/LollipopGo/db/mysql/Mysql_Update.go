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
func (this *mysql_db) Modefy_PlayerDataGM(uid, itype, number int) {
	// 默认修改VIP
	strSql := "update userinfo_copy set vip=? where uid=?"
	if itype == Proto2.MODIFY_COIN {
		strSql = "update userinfo_copy set vip=? where uid=?"
	} else if itype == Proto2.MODIFY_LEV {
		strSql = "update userinfo_copy set vip=? where uid=?"
	} else if itype == Proto2.MODIFY_LEV {
		strSql = "update userinfo_copy set vip=? where uid=?"
	} else if itype == Proto2.MODIFY_LEV {
		strSql = "update userinfo_copy set vip=? where uid=?"
	}
	// 修改数据
	stmt, err := this.STdb.Prepare(strSql)
	CheckErr(err)
	res, err := stmt.Exec(number, uid)
	affect, err := res.RowsAffected()
	fmt.Println("更新数据：", affect)
	CheckErr(err)
	return
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
