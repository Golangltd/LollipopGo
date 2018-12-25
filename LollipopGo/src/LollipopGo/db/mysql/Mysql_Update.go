package Mysyl_DB

import (
	"database/sql"
	"fmt"
)

/*
   更新数据库数据操作
*/

// GM 或者数据更新玩家数据操作
// 根据GM的命令的数据的类型来处理
func Update_PlayerData() {
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
