/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月3日
*/

package config

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var G_StInfoBaseST map[string]*DBBaseConfig

type DBBaseConfig struct {
	ID        string
	LoginName string // 数据库的登录名
	LoginPW   string // 数据库的登录密码
	DBIP      string // 数据库的IP
	DBPort    string // 数据库的端口（默认3306）
	Type      string // 数据库的类型
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
		Infotmp.Type = ss[i][5]
		G_StInfoBaseST[Infotmp.ID] = Infotmp
	}
	fmt.Println(G_StInfoBaseST)
	return true
}
