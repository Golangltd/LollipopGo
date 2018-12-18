package main

import (
	"Proto"
	"Proto/Proto2"
	"encoding/json"
	"fmt"
	"glog-master"
	"go-concurrentMap-master"
	"reflect"

	"code.google.com/p/go.net/websocket"
)

// 网络数据结构的保存
type NetDataConn struct {
	Connection    *websocket.Conn           //  链接信息
	StrMd5        string                    //  校验加密的信息
	MapSafe       *concurrent.ConcurrentMap //  保存client的信息
	MapSafeServer *concurrent.ConcurrentMap //  保存server的信息
}

// 结构体数据类型
type Requestbody struct {
	req string
}

//json转化为map:数据的处理
func (r *Requestbody) Json2map() (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
		glog.Error("Json2map:", err.Error())
		return nil, err
	}
	return result, nil
}

// func (r *Requestbody) Json2map() (s map[interface{}]interface{}, err error) {
// 	var result map[interface{}]interface{}
// 	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
// 		glog.Error("Json2map:", err.Error())
// 		return nil, err
// 	}
// 	return result, nil
// }

func (this *NetDataConn) PullFromClient() {
	for {

		var content string
		if err := websocket.Message.Receive(this.Connection, &content); err != nil {
			break
		}
		if len(content) == 0 || len(content) >= 4096 {
			break
		}
		go this.SyncMeassgeFun(content)
	}
}

func (this *NetDataConn) SyncMeassgeFun(content string) {
	var r Requestbody
	r.req = content

	if ProtocolData, err := r.Json2map(); err == nil {
		// 处理我们的函数
		this.HandleCltProtocol(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData)
	} else {
		glog.Error("解析失败：", err.Error())
	}
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// 处理函数(底层函数了，必须面向所有的数据处理)
//func (this *NetDataConn) HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[interface{}]interface{}) {
func (this *NetDataConn) HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}) {

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
			this.PlayerSendMessage(ErrorST)
		}
	}()

	// 分发处理  --- 首先判断主协议存在，再判断子协议存在不

	//glog.Info(protocol)
	//glog.Info(Proto.GameData_Proto)

	//类型
	glog.Info(typeof(protocol))
	glog.Info(typeof(protocol2))
	//glog.Info(typeof(Proto.GameData_Proto))

	switch protocol {
	case float64(Proto.G_GateWay_Proto):
		{
			// 网关协议
			this.HandleCltProtocol2GW(protocol2, ProtocolData)
		}
	case float64(Proto.GameData_Proto):
		{
			// 子协议处理
			this.HandleCltProtocol2(protocol2, ProtocolData)

		}
	case float64(Proto.GameDataDB_Proto):
		{ // DB_server

		}
	case float64(Proto.G_GameGlobal_Proto):
		{ // global_server
			this.HandleCltProtocol2GL(protocol2, ProtocolData)
		}
	case float64(Proto.GameNet_Proto):
		{
			this.HandleCltProtocol2Net(protocol2, ProtocolData)
		}
	case float64(Proto.G_Snake_Proto):
		{ // 贪吃蛇的主协议
			fmt.Println("贪吃蛇的主协议!!!")
			this.HandleCltProtocol2Snake(protocol2, ProtocolData)

		}
	default:
		panic("主协议：不存在！！！")
	}
	return
}

// 子协议的处理
func (this *NetDataConn) HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.C2S_PlayerLoginProto2):
		{
			// 功能函数处理 --  用户登陆协议
			this.PlayerLogin(ProtocolData)
		}
	case float64(Proto2.C2S_PlayerRunProto2):
		{
			// 功能函数处理 --  用户行走、奔跑
			this.PlayerRun(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

// 用户奔跑的协议
func (this *NetDataConn) PlayerRun(ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil {
		panic(" 主协议 GameData_Proto ，子协议 C2S_PlayerRunProto2,玩家行走功能数据错误！")
		return
	}

	StrOpenID := ProtocolData["OpenID"].(string)
	StrRunX := ProtocolData["StrRunX"].(string)
	StrRunY := ProtocolData["StrRunY"].(string)
	StrRunZ := ProtocolData["StrRunZ"].(string)

	// 广播协议
	data := &Proto2.S2C_PlayerRun{
		Protocol:  Proto.GameData_Proto,
		Protocol2: Proto2.S2C_PlayerRunProto2,
		OpenID:    StrOpenID,
		StrRunX:   StrRunX,
		StrRunY:   StrRunY,
		StrRunZ:   StrRunZ,
	}
	// 发送数据给客户端了
	//Broadcast(data)
	this.PlayerSendMessage(data)
	return
}

// 用户登陆的协议
func (this *NetDataConn) PlayerLogin(ProtocolData map[string]interface{}) {
	// 服务器的逻辑处理
	// 获取用户发过来的数据的信息
	if ProtocolData["StrLoginName"] == nil ||
		ProtocolData["StrLoginPW"] == nil ||
		ProtocolData["StrLoginEmail"] == nil {
		panic(" 主协议 GameData_Proto ，子协议 C2S_PlayerLoginProto2,登陆功能数据错误！")
		return
	}
	// 玩家的登陆名字
	StrLoginName := ProtocolData["StrLoginName"].(string)
	StrLoginPW := ProtocolData["StrLoginPW"].(string)
	StrLoginEmail := ProtocolData["StrLoginEmail"].(string)

	glog.Info(StrLoginName, StrLoginPW, StrLoginEmail)
	// 数据库的保存 -- 发给DBserver
	// 放给我们的 客户端
	// channel 操作
	// 保存玩家数据
	playerdata := &NetDataConn{
		Connection: this.Connection,
		StrMd5:     (StrLoginName + StrLoginPW),
		MapSafe:    this.MapSafe,
	}
	// 保存 --
	// 优化： 讲并发非安全的-->并发安全的数据结构
	//	glog.Info("-------------------------------")
	this.MapSafe.Put(StrLoginName+"PlayerUID"+"|connect", playerdata)
	//	glog.Info(this.MapSafe)
	glog.Info("-------------------------------")

	glog.Info(G_PlayerData["123456"])
	// 服务器-->客户端
	data := &Proto2.S2C_PlayerLogin{
		Protocol:   Proto.GameData_Proto,
		Protocol2:  Proto2.S2C_PlayerLoginProto2,
		PlayerData: nil,
	}
	// 发送数据给客户端了
	// Broadcast(data)
	this.PlayerSendMessage(data)
	return
}

// 广播函数处理
func Broadcast(data interface{}) {

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
						// 发送数据
						v.(interface{}).(*NetDataConn).PlayerSendMessage(data)
					}
				}
			}
		}
	}

	return
}
