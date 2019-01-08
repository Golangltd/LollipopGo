package main

import (
	_ "LollipopGo/db/mysql"
	_ "LollipopGo/db/redis"
	"cache2go"
	"encoding/base64"
	"flag"
	"fmt"
	"go-concurrentMap-master"
	"strings"

	"code.google.com/p/go.net/websocket"
)

// 全局的网络信息结构
var G_PlayerData map[string]*NetDataConn
var G_PlayerNet map[string]int // 心跳结构信息存储的结构
var G_PlayerNetSys map[string]int
var G_Net_Count map[string]int
var M *concurrent.ConcurrentMap       // 并发安全的client的链接
var MRoom *concurrent.ConcurrentMap   // 并发安全的房间的链接
var MServer *concurrent.ConcurrentMap // 并发安全的server链接
var addr = flag.String("addr", "127.0.0.1:8888", "http service address")
var WS *websocket.Conn
var icount, icounttmp int
var cacheGW *cache2go.CacheTable // 网关cache
var DSQGameID = 10001

// server data;推送数据时候用
var strGlobalServer string = ""
var strDSQServer string = ""

// 游戏服务器的初始化
func init() {
	// 初始化
	G_PlayerData = make(map[string]*NetDataConn)
	G_PlayerNet = make(map[string]int)
	G_PlayerNetSys = make(map[string]int)
	G_Net_Count = make(map[string]int)
	// 并发安全的初始化
	M = concurrent.NewConcurrentMap()
	MRoom = concurrent.NewConcurrentMap()
	MServer = concurrent.NewConcurrentMap()
	cacheGW = cache2go.Cache("LollipopGo_GateWay")
	// go G_timer()
	// go G_timeout_kick_Player()
	// redis 测试
	// go Redis_DB.INIT()
	return
}

func Go_func() {
	fmt.Println("Golang语言社区")
	return
}

// 服务器的初始化
func GameServerINIT() {
	// 1 主动链接网关服  -- 获取IP+端口
	// 2 建立与网关服的心跳-- 确保server正常 -- kill
	url := "ws://" + *addr + "/GolangLtd"
	WS, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	go GS2GW_Timer(WS)
	go GameServerReceive(WS)
}

// 处理数据的返回
func GameServerReceive(ws *websocket.Conn) {
	for {
		var content string
		err := websocket.Message.Receive(ws, &content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// decode
		fmt.Println(strings.Trim("", "\""))
		fmt.Println(content)
		content = strings.Replace(content, "\"", "", -1)
		contentstr, errr := base64Decode([]byte(content))
		if errr != nil {
			fmt.Println(errr)
		}
		icounttmp++
		fmt.Println("返回数据：", string(contentstr))
	}
}

// 解码
func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
