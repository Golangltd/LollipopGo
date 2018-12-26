package conf

import (
	"LollipopGo/LollipopGo/log"
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

// 加载服务器配置
func init() {
	// 基础配置
	if true {
		data, err := ioutil.ReadFile("conf/server.json")
		if err != nil {
			log.Debug("-------------%v", err)
		}
		err = json.Unmarshal(data, &Server)
		if err != nil {
			log.Debug("+++++++++++++%v", err)
		}
	}
	// 服务器配置
	if true {
		data, err := ioutil.ReadFile("conf/cluster.json")
		if err != nil {
			log.Debug("-------------%v", err)
		}
		err = json.Unmarshal(data, &ServerConf)
		if err != nil {
			log.Debug("+++++++++++++%v", err)
		}
	}

}
