/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月3日
*/

package config

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	//"strconv"
	"strings"
)

var G_StInfoBaseST map[string]*DBBaseConfig

type DBBaseConfig struct {
	ID        string
	LoginName string // 数据库的登录名
	LoginPW   string // 数据库的登录密码
	DBIP      string // 数据库的IP
	DBPort    string // 数据库的端口（默认3306）
	DBName    string // 数据库名字
	Itype     string // 数据库类型
}

//获取配置信息
func init() {
	ReadCsv_ConfigFile_StCard2List_Fun()
	return
}

// 获取配置信息
func ReadCsv_ConfigFile_StCard2List_Fun() bool {
	// 获取数据，按照文件
	fileName := "config.csv"
	fileName = "./" + fileName
	cntb, err := ioutil.ReadFile(fileName)
	if err != nil {
		return false
	}
	// 读取文件数据
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	ss, _ := r2.ReadAll()
	sz := len(ss)

	// 循环取数据
	for i := 1; i < sz; i++ {

		Infotmp := new(DBBaseConfig)
		Infotmp.ID = ss[i][0]
		Infotmp.LoginName = ss[i][1]
		Infotmp.LoginPW = ss[i][2]
		Infotmp.DBIP = ss[i][3]
		Infotmp.DBPort = ss[i][4]
		Infotmp.DBName = ss[i][5]
		Infotmp.Itype = ss[i][6]
		G_StInfoBaseST[Infotmp.ID] = Infotmp
	}
	fmt.Println(G_StInfoBaseST)
	return true
}

// 链接池的最大链接数量
const MAX_POOL_SIZE int = 200

// 全局数据库变量
var MySQLPool chan *sql.DB

// 获取数据链接
func getMySQL() *sql.DB {
	// 获取链接
	conn := GetMySQL1()
	// 压入队列
	putMySQL(conn)
	return conn
}

// 获取数据链接
func GetMySQL() *sql.DB {
	// 获取链接
	conn := GetMySQL1()
	// 压入队列
	putMySQL(conn)
	return conn
}

// 获取链接指针函数
func GetMySQL1() *sql.DB {

	if MySQLPool == nil {
		MySQLPool = make(chan *sql.DB, MAX_POOL_SIZE)
	}
	if len(MySQLPool) == 0 {
		go func() {
			for i := 0; i < MAX_POOL_SIZE/2; i++ {
				mysql := new(sql.DB)
				var err error
				var StrConnection = ""
				//if Log_Eio.BTest == true {
				// 测试
				StrConnection = "root" + ":" + "123456" + "@tcp(" + "127.0.0.1" + ":3306)/" + "gl_test"
				//}
				mysql, err = sql.Open("mysql", StrConnection)
				if err != nil {

					continue
				}
				putMySQL(mysql)
			}
		}()
	}
	return <-MySQLPool
}

//存储指针函数
func putMySQL(conn *sql.DB) {
	if MySQLPool == nil {
		MySQLPool = make(chan *sql.DB, MAX_POOL_SIZE)
	}
	if len(MySQLPool) == MAX_POOL_SIZE {
		conn.Close()
		return
	}
	MySQLPool <- conn
}
