package conf

import (
	"encoding/json"
	"io/ioutil"
)

// 服务器集群配置
var ServerConf struct {
	LoginServerAddr  string
	GateWayAddr      string
	DBServerAddr     string
	GlobalServerAddr string
}

// 服务器结构
var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
}

// 数据库配置
var DBServer struct {
	MasterLoginName     string
	MasterLoginPassword string
	SlaveLoginName      string
	SlaveLoginPassword  string
	MaxOpenConns        string
	MaxIdleConns        string
	MasterMysql_IP      string
	MasterMysql_Port    string
	SlaveMysql_IP       string
	SlaveMysql_Port     string
}

func init() {
	// 基础配置
	if true {
		data, _ := ioutil.ReadFile("conf/server.json")
		json.Unmarshal(data, &Server)

	}

	// 服务器配置
	if true {
		data, _ := ioutil.ReadFile("conf/cluster.json")
		json.Unmarshal(data, &ServerConf)

	}

	//  读取数据库配置
	if true {
		data, _ := ioutil.ReadFile("conf/mysql.json")
		json.Unmarshal(data, &DBServer)

	}

}
