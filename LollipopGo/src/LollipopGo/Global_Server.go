package main

import (
	"LollipopGo/LollipopGo/log"
	"Proto"
	"Proto/Proto2"
	"flag"
	"fmt"
	"strings"

	"LollipopGo/LollipopGo/util"

	"code.google.com/p/go.net/websocket"
)

/*
  匹配、活动服务器
	1 匹配玩家活动
*/

var addrG = flag.String("addrG", "127.0.0.1:8888", "http service address")
var Conn *websocket.Conn // 保存用户的链接信息，数据会在主动匹配成功后进行链接

// 初始化操作
func init() {
	if initGateWayNet() {
		panic("链接 gateway server 失败!")
		return
	}
	return
}

func initGateWayNet() bool {

	return false

	fmt.Println("用户客户端客户端模拟！")
	log.Debug("用户客户端客户端模拟！")
	url := "ws://" + *addrG + "/GolangLtd"
	conn, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return false
	}
	Conn = conn
	// 协程支持  --接受线程操作 全球协议操作
	go GameServerReceiveG(Conn)
	// 发送链接的协议 ---》
	initConn(Conn)
	return true
}

// 处理数据
func GameServerReceiveG(ws *websocket.Conn) {
	for {
		var content string
		err := websocket.Message.Receive(ws, &content)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// decode
		fmt.Println(strings.Trim("", "\""))
		fmt.Println(content)
		content = strings.Replace(content, "\"", "", -1)
		contentstr, errr := base64Decode([]byte(content))
		if errr != nil {
			fmt.Println(errr)
			continue
		}
		// 解析数据 --
		fmt.Println("返回数据：", string(contentstr))
		go SyncMeassgeFunG(string(contentstr))
	}
}

// 链接分发 处理
func SyncMeassgeFunG(content string) {
	var r Requestbody
	r.req = content

	if ProtocolData, err := r.Json2map(); err == nil {
		// 处理我们的函数
		HandleCltProtocolG(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData)
	} else {
		log.Debug("解析失败：", err.Error())
	}
}

//  主协议处理
func HandleCltProtocolG(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			//发消息给客户端
			ErrorST := Proto2.G_Error_All{
				Protocol:  Proto.G_Error_Proto,      // 主协议
				Protocol2: Proto2.G_Error_All_Proto, // 子协议
				ErrCode:   "80006",
				ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
			}
			// 发送给玩家数据
			fmt.Println("贪吃蛇的主协议!!!", ErrorST)
		}
	}()

	// 协议处理
	switch protocol {
	case float64(Proto.G_Snake_Proto):
		{ // Global Server 主要协议处理
			fmt.Println("Global server 主协议!!!")
			HandleCltProtocol2Glogbal(protocol2, ProtocolData)

		}
	default:
		panic("主协议：不存在！！！")
	}
	return
}

// 子协议的处理
func HandleCltProtocol2Glogbal(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.GW2G_ConnServerProto2):
		{ // 网关返回数据
			fmt.Println("gateway server 返回给global server 数据信息！！！")

		}
	case float64(Proto2.G2GW_PlayerEntryHallProto2):
		{ // 网关请求获取大厅数据
			fmt.Println("玩家请求获取大厅数：默认获奖列表、跑马灯等")
			G2GW_PlayerEntryHallProto2Fucn(Conn)
		}

	default:
		panic("子协议：不存在！！！")
	}
	return
}

// 返回给玩家数据
func G2GW_PlayerEntryHallProto2Fucn(conn *websocket.Conn) {
	// 返回数据给GateWay

	iGamePlayerNum := make(map[string]int)
	iGamePlayerNum["1001"] = 1000
	iGamePlayerNum["1002"] = 9999

	// 组装数据
	data := &Proto2.GW2G_PlayerEntryHall{
		Protocol:      Proto.G_GameGlobal_Proto, // 游戏主要协议
		Protocol2:     Proto2.GW2G_PlayerEntryHallProto2,
		GamePlayerNum: iGamePlayerNum,
		DefaultAward:  nil,
		DefaultMsg:    nil,
	}
	fmt.Println(data)
	// 2 发送数据到服务器
	PlayerSendToServer(conn, data)
	return

}

// 链接到网关
func initConn(conn *websocket.Conn) {
	// 组装数据
	data := &Proto2.G2GW_ConnServer{
		Protocol:  Proto.G_GameGlobal_Proto, // 游戏主要协议
		Protocol2: Proto2.G2GW_ConnServerProto2,
		ServerID:  util.MD5_LollipopGO("8894" + "Global server"),
	}
	// fmt.Println(data)
	// 2 发送数据到服务器
	PlayerSendToServer(conn, data)
	return
}
