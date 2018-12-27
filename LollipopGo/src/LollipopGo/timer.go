package main

import (
	"Proto"
	"Proto/Proto2"
	"glog-master"
	"os"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
)

// 字符串分割函数
func Strings_Split(Data string, Split string) []string {
	return strings.Split(Data, Split)
}

// 超时踢人
func G_timeout_kick_Player() {
	for {
		select {
		case <-time.After(10 * time.Second):
			{
				// 1 获取我们心跳数据--玩家的 测试一个玩家 data[A] = 1
				// 2 玩家的心跳保存下来-- 临时的保存  datatmp[A] = 1
				// 3 每10s 对比一次：临时的保存的数据 与我们心跳的数据是否相同 data[A] == datatmp[A]
				// 4 30s 还是没有变化  kick player  A
				// 并发安全map优化：
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

// 数据推送给客户端定时
func G_timer() {
	for {
		select {
		case <-time.After(20 * time.Second):
			{
				if len(G_PlayerData) == 0 {
					continue
				}

				if G_PlayerData["123456"] != nil {
					data := &Proto2.S2C_PlayerLogin{
						Protocol:   Proto.GameData_Proto,
						Protocol2:  Proto2.S2C_PlayerLoginProto2,
						PlayerData: nil,
					}
					G_PlayerData["123456"].PlayerSendMessage(data)
					glog.Info("发送数据：", data)
				}
			}
		}
	}
	return
}

// 游戏服务器---网关之间的通信
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
