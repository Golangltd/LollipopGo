package Mysyl_DB

import (
	"database/sql"
	"fmt"
)

/*
   删除数据库数据操作
*/

func DeleteFromDB(db *sql.DB, autid int) {
	stmt, err := db.Prepare("delete from userinfo where autid=?")
	CheckErr(err)
	res, err := stmt.Exec(autid)
	affect, err := res.RowsAffected()
	fmt.Println("删除数据：", affect)
}
