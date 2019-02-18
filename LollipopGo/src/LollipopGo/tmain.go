package main

import (
	"LollipopGo/LollipopGo"
	LollipopGoconf "LollipopGo/LollipopGo/conf"
	_ "LollipopGo/LollipopGo/match"
	"LollipopGo/conf"
	"Proto"
	"Proto/Proto2"
	"glog-master"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
)

var strServerType = "GW"
var GL_type = "G"

func init() {
	LollipopGoconf.LogLevel = conf.Server.LogLevel
	LollipopGoconf.LogPath = conf.Server.LogPath
	LollipopGoconf.LogFlag = conf.LogFlag
	LollipopGoconf.ConsolePort = conf.Server.ConsolePort
	LollipopGoconf.ProfilePath = conf.Server.ProfilePath
	LollipopGo.Run()
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
	strServerType_GW := "GW"
	strServerType_GS := "GS"
	strServerType_DB := "DB"
	strServerType_DT := "DT"
	strServerType_GM := "GM"
	strServerType_GL := "GL"
	strServerType_Snake := "Snake"
	strServerType_DSQ := "DSQ"
	if len(os.Args) > 1 {
		strport = os.Args[1]
		strServerType = os.Args[2]
	}

	glog.Info(strport)
	glog.Info(strServerType)
	glog.Info(strServerType_GW)

	if "GW" == strServerType {
		glog.Info("Golang语言社区  gw")
		strServerType_GW = strServerType
	}
	glog.Info("Golang语言社区")
	glog.Flush()
	if strServerType == strServerType_GW {
		http.Handle("/GolangLtd", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	} else if strServerType == strServerType_GS {
		strport = "8889"
		go GameServerINIT()
		http.Handle("/GolangLtdGS", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	} else if strServerType == strServerType_DB {
		strport = "8890"
		MainListener(strport)
	} else if strServerType == strServerType_DT {
		strport = "8891"
		http.HandleFunc("/GolangLtdDT", IndexHandler)
		http.ListenAndServe(":"+strport, nil)
	} else if strServerType == strServerType_GM {
		strport = "8892"
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
		http.Handle("/GolangLtdGL", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	} else if strServerType == strServerType_DSQ {
		strport = "8885"
		// go GameServerINIT()
		http.Handle("/GolangLtdDSQ", websocket.Handler(wwwGolangLtd))
		if err := http.ListenAndServe(":"+strport, nil); err != nil {
			glog.Error("网络错误", err)
			return
		}
	}
	panic("【服务器类型】不存在")
}

func Strings_Split(Data string, Split string) []string {
	return strings.Split(Data, Split)
}

func G_timeout_kick_Player() {
	for {
		select {
		case <-time.After(10 * time.Second):
			{
				for itr := M.Iterator(); itr.HasNext(); {
					k, v, _ := itr.Next()
					// 取分隔符
					strsplit := Strings_Split(k.(string), "|")
					for i := 0; i < len(strsplit); i++ {
						if len(strsplit) < 2 {
							continue
						}
						// 进行数据的查询类型
						switch v.(interface{}).(type) {
						case *NetDataConn:
							{
								// 判断 链接是不是 connect
								if "" == "connect" {
									data := &Proto2.Net_Kicking_Player{
										Protocol:  Proto.GameNet_Proto,
										Protocol2: Proto2.Net_Kicking_PlayerProto2,
										ErrorCode: 10001,
									}
									// 发送数据
									v.(interface{}).(*NetDataConn).PlayerSendMessage(data)
								}
							}
						}
					}
				}
				// -------------------------------------------------------------

				if G_Net_Count["12345"] >= 3 {
					// 踢人
					data := &Proto2.Net_Kicking_Player{
						Protocol:  Proto.GameNet_Proto,
						Protocol2: Proto2.Net_Kicking_PlayerProto2,
						ErrorCode: 10001,
					}
					G_PlayerData["123456"].PlayerSendMessage(data)
					// 关闭链接
					G_PlayerData["123456"].Connection.Close()
					G_Net_Count["12345"] = 0
					continue
				}

				if len(G_PlayerNetSys) == 0 {
					G_PlayerNetSys["12345"] = G_PlayerNet["12345"]
				} else {
					if G_PlayerNetSys["12345"] == G_PlayerNet["12345"] {
						G_Net_Count["12345"]++
					}
				}
			}
		}
	}

}

func GS2GW_Timer(ws *websocket.Conn) {
	for {
		select {
		case <-time.After(5 * time.Second):
			{
				// 1 组装
				data := &Proto2.Net_HeartBeat{
					Protocol:  Proto.GameNet_Proto,
					Protocol2: Proto2.Net_HeartBeatProto2,
					OpenID:    "12345123451234512345123451234512345123451234512345123451234512345",
				}
				// 3 发送数据到服务器
				if ws != nil {
					PlayerSendToServer(ws, data)
					glog.Info("发送数据----：", data)
					icount++
					if icounttmp == icount-10 {
						os.Exit(0)
					}
					continue
				}
				glog.Info("发送数据：", data)

			}
		}
	}
	return
}
