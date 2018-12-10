package Mysyl_DB

import (
	"database/sql"
	"fmt"
)

/*
   插入数据库数据操作
*/

func insertToDB(db *sql.DB) {
	uid := GetNowtimeMD5()
	nowTimeStr := GetTime()
	stmt, err := db.Prepare("insert userinfo set username=?,departname=?,created=?,password=?,uid=?")
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
