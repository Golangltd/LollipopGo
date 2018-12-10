package Mysyl_DB

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db = &mysql_db{}
	db.mysql_open()
	return
}

func (f *mysql_db) mysql_open() {
	Odb, err := sql.Open("mysql", dbusername+":"+dbpassowrd+"@tcp("+dbhostsip+")/"+dbname)
	if err != nil {
		fmt.Println("链接失败")
	}
	fmt.Println("链接数据库成功...........已经打开")
	// 设置链接池
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)
	db.Ping()
	f.db = Odb
}
