package Mysyl_DB

import (
	LollipopGodb "LollipopGo/LollipopGo/db"
	"LollipopGo/conf"
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//--------------------------------------------------------------------------
	// 加载配置
	LollipopGodb.MasterLoginName = conf.DBServer.MasterLoginName
	LollipopGodb.MasterLoginPassword = conf.DBServer.MasterLoginPassword
	LollipopGodb.MaxOpenConns = conf.DBServer.MaxOpenConns
	LollipopGodb.MaxIdleConns = conf.DBServer.MaxIdleConns
	LollipopGodb.MasterMysql_IP = conf.DBServer.MasterMysql_IP
	LollipopGodb.MasterMysql_Port = conf.DBServer.MasterMysql_Port
	//--------------------------------------------------------------------------
	DB = &mysql_db{}
	DB.mysql_open()
	return
}

func (this *mysql_db) mysql_open() {
	Odb, err := sql.Open("mysql", dbusername+":"+dbpassowrd+"@tcp("+dbhostsip+")/"+dbname)
	if err != nil {
		fmt.Println("链接失败")
	}
	fmt.Println("链接数据库成功...........已经打开")
	this.STdb = Odb
	// 设置链接池
	this.STdb.SetMaxOpenConns(dbMaxOpenConns)
	this.STdb.SetMaxIdleConns(dbMaxIdleConns)
	this.STdb.Ping()
	defer Odb.Close()
}
