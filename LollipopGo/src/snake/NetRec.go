package main

import (
	"Proto"
	"Proto/Proto2"
	"encoding/json"
	"fmt"
	"strings"

	"LollipopGo/LollipopGo/log"

	"code.google.com/p/go.net/websocket"
)

// 处理数据的返回
func GameServerReceive(ws *websocket.Conn) {
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
		go SyncMeassgeFun(string(contentstr))
	}
}

// 结构体数据类型
type Requestbody struct {
	req string
}

// json转化为map:数据的处理
func (r *Requestbody) Json2map() (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
		log.Debug("Json2map:", err.Error())
		return nil, err
	}
	return result, nil
}

func SyncMeassgeFun(content string) {
	var r Requestbody
	r.req = content

	if ProtocolData, err := r.Json2map(); err == nil {
		// 处理我们的函数
		HandleCltProtocol(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData)
	} else {
		log.Debug("解析失败：", err.Error())
	}
}

// 字符串 解析成 json

func HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}) {
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
		{ // 贪吃蛇的主协议
			fmt.Println("贪吃蛇的主协议!!!")
			HandleCltProtocol2Snake(protocol2, ProtocolData)

		}
	default:
		panic("主协议：不存在！！！")
	}
	return
}

// 子协议的处理
func HandleCltProtocol2Snake(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.S2S_PlayerLoginSProto2):
		{
			fmt.Println("贪吃蛇:玩家进入游戏的协议!!!")
			// 玩家进入游戏的协议
			EntryGameSnake(ProtocolData)
		}
	case float64(Proto2.S2S_PlayerEntryGameProto2):
		{
			fmt.Println("贪吃蛇:玩家匹配成功协议!!!")
			// 玩家进入游戏的协议
			// EntryGameSnake(ProtocolData)
		}

	default:
		panic("子协议：不存在！！！")
	}

	return
}

func EntryGameSnake(ProtocolData map[string]interface{}) {
	StrToken := ProtocolData["Token"].(string)
	fmt.Println("贪吃蛇:玩家进入游戏的协议!!!", StrToken)
	return
}
