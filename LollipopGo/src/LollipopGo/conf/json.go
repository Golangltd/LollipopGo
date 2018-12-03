package conf

import (
	"LollipopGo/LollipopGo/log"
	"encoding/json"
	"io/ioutil"
)

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
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		log.Debug("-------------%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Debug("+++++++++++++%v", err)
	}
}
