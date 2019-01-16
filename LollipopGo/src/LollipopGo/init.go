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

var G_PlayerData map[string]*NetDataConn
var G_PlayerNet map[string]int
var G_PlayerNetSys map[string]int
var G_Net_Count map[string]int
var M *concurrent.ConcurrentMap
var MRoom *concurrent.ConcurrentMap
var MServer *concurrent.ConcurrentMap
var addr = flag.String("addr", "127.0.0.1:8888", "http service address")
var WS *websocket.Conn
var icount, icounttmp int
var cacheGW *cache2go.CacheTable
var DSQGameID = 10001

var strGlobalServer string = ""
var strDSQServer string = ""

func init() {

	G_PlayerData = make(map[string]*NetDataConn)
	G_PlayerNet = make(map[string]int)
	G_PlayerNetSys = make(map[string]int)
	G_Net_Count = make(map[string]int)

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

func GameServerINIT() {
	url := "ws://" + *addr + "/GolangLtd"
	WS, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	go GS2GW_Timer(WS)
	go GameServerReceive(WS)
}

func GameServerReceive(ws *websocket.Conn) {
	for {
		var content string
		err := websocket.Message.Receive(ws, &content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
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

func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
