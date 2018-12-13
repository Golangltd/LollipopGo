package main

import (
	"LollipopGo/LollipopGo"
	LollipopGoconf "LollipopGo/LollipopGo/conf"
	"LollipopGo/conf" // 服务器配置——针对不同环境的配置
	"glog-master"
	"net/http"
	_ "net/http/pprof"
	"os"

	"code.google.com/p/go.net/websocket"
)

func init() {
	// 加载配置
	LollipopGoconf.LogLevel = conf.Server.LogLevel
	LollipopGoconf.LogPath = conf.Server.LogPath
	LollipopGoconf.LogFlag = conf.LogFlag
	LollipopGoconf.ConsolePort = conf.Server.ConsolePort
	LollipopGoconf.ProfilePath = conf.Server.ProfilePath
	glog.Info("LollipopGo version:", LollipopGo.version)
}

func main() {
	// os.Args[0] == 执行文件的名字
	// os.Args[1] == 第一个参数
	// os.Args[2] == 类型 Client -websocket-> GW -websocket/rpc-> GS -websocket/rpc-> DB
	glog.Info(os.Args[:])
	if len(os.Args[:]) < 3 {
		panic("参数小于2个！！！ 例如：xxx.exe +【端口】+【服务器类型】")
		return
	}
	strport := "8888"
	strServerType := "GW"
	strServerType_GW := "GW"
	strServerType_GS := "GS"
	strServerType_DB := "DB"
	strServerType_DT := "DT"
	strServerType_GM := "GM"
	strServerType_GL := "GL"
	strServerType_Snake := "Snake"
	if len(os.Args) > 1 {
		strport = os.Args[1]
		strServerType = os.Args[2]
	}
	glog.Info(strport)
	glog.Info("Golang语言社区")
	glog.Flush()
	if strServerType == strServerType_GW {
		http.Handle("/GolangLtd", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	} else if strServerType == strServerType_GS {
		strport = "8889"    // 多个 -- server
		go GameServerINIT() // 游戏服务器的初始化操作
		http.Handle("/GolangLtdGS", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	} else if strServerType == strServerType_DB {
		strport = "8890"
		MainListener(strport)
		// http.HandleFunc("/GolangLtdDT", GolangLtdDB)
		// http.ListenAndServe(":"+strport, nil)
	} else if strServerType == strServerType_DT {
		strport = "8891" //  登录服务器 -- 大厅服务器
		http.HandleFunc("/GolangLtdDT", IndexHandler)
		http.ListenAndServe(":"+strport, nil)
	} else if strServerType == strServerType_GM {
		strport = "8892" //  GM 系统操作 -- 修改金币等操作
		http.HandleFunc("/GolangLtdGM", IndexHandlerGM)
		http.ListenAndServe(":"+strport, nil)
	} else if strServerType == strServerType_Snake {
		strport = "8893"
		http.Handle("/GolangLtdSnake", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	} else if strServerType == strServerType_GL {
		strport = "8894"
		http.Handle("/GolangLtdGL", websocket.Handler(wwwGolangLtdGL))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	}
	panic("【服务器类型】不存在")
}
