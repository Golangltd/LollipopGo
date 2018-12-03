package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"Proto"
	"Proto/Proto2"
	"encoding/base64"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"code.google.com/p/go.net/websocket"
)

var GW sync.WaitGroup
var Scount int

func main1() {

	fmt.Println(os.Args[1:])
	fmt.Println(flag.Arg(1))
	fmt.Println("Entry server count:", os.Args[1]) // 人数

	t1 := time.Now()
	count := os.Args[1]
	loops, _ := strconv.Atoi(os.Args[2]) // 并发次数
	int1, _ := strconv.Atoi(count)
	GW.Add(int1)
	Scount = int1 * loops
	for i := 0; i < int1; i++ {
		go GoroutineFunc(loops)
	}
	GW.Wait()
	elapsed := time.Since(t1)
	fmt.Println("Total count:", int1*loops)
	fmt.Println("Success count:", Scount)
	fmt.Println("Cysle TPS:", float64(int1*loops)/elapsed.Seconds())
	fmt.Println("Taken Time(s) :", elapsed)
	fmt.Println("Average Latency time(ms):", elapsed.Seconds()*1000/(float64(int1*loops)))
	//-------------------------------------------------------------------

}
func GoroutineFunc(loops int) {
	fmt.Println("Robot 客户端模拟！")
	url := "ws://" + *addr + "/GolangLtd"
	ws, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	// 数据的发送

	for i := 0; i < loops; i++ {
		Send(ws, "HeartBeat")
		// 发送心跳的协议
		// 数据处理
		var content string
		err := websocket.Message.Receive(ws, &content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// decode
		//fmt.Println(strings.Trim("", "\""))
		//fmt.Println(content)
		content = strings.Replace(content, "\"", "", -1)
		contentstr, errr := base64Decode([]byte(content))
		if errr != nil {
			fmt.Println(errr)
		}
		// 解析数据
		//fmt.Println(string(contentstr))
		_ = contentstr
	}
	GW.Done()

}

/*
robot
1 模拟玩家的正常的“操作”，例如 行走 跳跃 开枪等等
2 做服务器的性能的测试，例如 并发量  内存 CPU 等等
3 压力测试

注意点：
1  模拟 ---> 多线程模拟  goroutine  --- server ！！！

首先：
1 net 网络使用websocket 进行连接
2  send  如何发送 ？？
*/
var addr = flag.String("addr", "127.0.0.1:8888", "http service address")
var connbak *websocket.Conn

// 1 robot客户端 我们是可以一起链接的 ---> websocket.Dial 每次都返回一个
// 2 多个 websocket.Dial   --->  多个客户端的链接

func main() {
	fmt.Println("Robot 客户端模拟！")
	url := "ws://" + *addr + "/GolangLtd"
	ws, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	// 数据的发送
	for i := 0; i < 100; i++ {
		// go Send(ws, "Login")
	}

	go Send(ws, "HeartBeat")
	// 发送心跳的协议
	connbak = new(websocket.Conn)
	connbak = ws
	go Timer(ws)

	// 数据处理
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
		// 解析数据
		fmt.Println(string(contentstr))
	}

}

// 解码
func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

// 消息的流程？
// 1 针对我们的消息结构进行数据的组装
// 2 针对我们组装的数据经行一个数据格式的转换 --> json
// 3 json 数据直接发送到我们server

// 发送函数
func Send(conn *websocket.Conn, Itype string) {

	if Itype == "Login" {

		// 1 组装
		data := &Proto2.C2S_PlayerLogin{
			Protocol:      Proto.GameData_Proto,
			Protocol2:     Proto2.C2S_PlayerLoginProto2,
			Itype:         1,
			StrLoginName:  "Golangltd" + Rand_LoginName(),
			StrLoginPW:    "1234556",
			StrLoginEmail: "123455@qq.com",
		}
		// 3 发送数据到服务器
		PlayerSendToServer(conn, data)
	} else if Itype == "HeartBeat" { // 心跳
		// 1 组装
		data := &Proto2.Net_HeartBeat{
			Protocol:  Proto.GameNet_Proto,
			Protocol2: Proto2.Net_HeartBeatProto2,
			OpenID:    "22323",
		}
		// 3 发送数据到服务器
		PlayerSendToServer(conn, data)
	}

	return
}

// 公用的send函数
func PlayerSendToServer(conn *websocket.Conn, data interface{}) {

	// 2 结构体转换成json数据
	jsons, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	///fmt.Println("jsons:", string(jsons))
	errq := websocket.Message.Send(conn, jsons)
	if errq != nil {
		fmt.Println(errq)
	}
	return
}

// 随机函数
func Rand_LoginName() string {
	idata := rand.Intn(100000)
	// 去重？
	strdata := strconv.Itoa(idata)
	return strdata
}
