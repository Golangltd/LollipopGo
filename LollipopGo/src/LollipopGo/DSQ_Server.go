package main

import (
	"LollipopGo/LollipopGo/log"
	"Proto"
	"Proto/Proto2"
	"flag"
	"fmt"
	_ "fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	"LollipopGo/LollipopGo/util"
	"LollipopGo/ReadCSV"

	"LollipopGo/LollipopGo/player"

	"code.google.com/p/go.net/websocket"
)

var addrDSQ = flag.String("addrDSQ", "127.0.0.1:8888", "http service address")
var ConnDSQ *websocket.Conn // 保存用户的链接信息，数据会在主动匹配成功后进行链接
var ConnDSQRPC *rpc.Client

// 初始化操作
func init() {
	if !initDSQGateWayNet() {
		fmt.Println("链接 gateway server 失败!")
		return
	}
	fmt.Println("链接 gateway server 成功!")
	// 初始化数据
	initDSQNetRPC()
	return
}

// 初始化RPC
func initDSQNetRPC() {
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		log.Debug("dial error:", err)
	}
	ConnDSQRPC = client
}

// 初始化网关
func initDSQGateWayNet() bool {

	fmt.Println("用户客户端客户端模拟！")
	log.Debug("用户客户端客户端模拟！")
	url := "ws://" + *addrDSQ + "/GolangLtd"
	conn, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return false
	}
	ConnDSQ = conn
	// 协程支持  --接受线程操作 全球协议操作
	go GameServerReceiveG(Conn)
	// 发送链接的协议 ---》
	initConnDSQ(Conn)
	return true
}

// 链接到网关
func initConnDSQ(conn *websocket.Conn) {
	// 协议修改
	data := &Proto2.G2GW_ConnServer{
		Protocol:  Proto.G_GameGlobal_Proto, // 游戏主要协议
		Protocol2: Proto2.G2GW_ConnServerProto2,
		ServerID:  util.MD5_LollipopGO("8895" + "DSQ server"),
	}
	fmt.Println(data)
	// 2 发送数据到服务器
	PlayerSendToServer(conn, data)
	return
}
