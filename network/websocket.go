package impl

import (
	"LollipopGo/Proxy_Server/Proto"
	"LollipopGo/global_Interface"
	"LollipopGo/util"
	"encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"github.com/json-iterator/go"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var BytebufLen int64 = 10000000
var IMsg MsgHandleClt.Msg_data

type OnlineUser struct {
	Connection *websocket.Conn
	inChan     chan string
	outChan    chan interface{}
	closeChan  chan int
	goExit     chan int
	isClosed   bool
	HandleClt  MsgHandleClt.Msg_data
}

func InitConnection(wsConn *websocket.Conn) (*OnlineUser, error) {
	conn := &OnlineUser{
		Connection: wsConn,
		inChan:     make(chan string, BytebufLen),
	}

	defer conn.Connection.Close()
	go conn.handleLoop()
	conn.readLoop()
	//go conn.readLoop()
	//select {}

	return conn, nil
}

// 20240710
func (this *OnlineUser) readLoop() {

	for {
		var content string
		err := websocket.Message.Receive(this.Connection, &content)
		if err != nil {
			//if err == io.EOF || err == io.ErrClosedPipe || content == "" || err == io.ErrNoProgress {
			if err == io.EOF {
				IMsg.CloseEOF(this.Connection)
				glog.Info("协程的数量 :", runtime.NumGoroutine())
				//this.Connection.Close()
				//runtime.Goexit()
				return
			}
			//break
			continue
		}
		select {
		case this.inChan <- content:
		case <-time.After(60 * time.Second):
			glog.Info("readLoop:超时----")
			//glog.Info("协程的数量 :", runtime.NumGoroutine())
			//this.Connection.Close()
			//runtime.Goexit()
			return
			//default:
			//	fmt.Println("Channel is empty, unable to read data")
		}
	}
}

func (this *OnlineUser) handleLoop() {

	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			glog.Info("异常捕获:", strerr)
		}
	}()

	for {
		if this.inChan == nil {
			continue
		}
		var r Requestbody
		select {
		case r.req = <-this.inChan:
		case <-time.After(200 * time.Second):
			glog.Info("handleLoop:超时----")
			//glog.Info("协程的数量 :", runtime.NumGoroutine())
			//this.Connection.Close()
			//runtime.Goexit()
			return
		}
		if len(r.req) <= 0 {
			continue
		}

		if ProtocolData, err := r.Json2map(); err == nil {
			IMsg.HandleCltProtocol(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData, this.Connection)
		} else {
			content := r.req
			content = strings.Replace(content, "\"", "", -1)
			contentstr, errr := base64Decode([]byte(content))
			if errr != nil {
				fmt.Println(errr)
				this.Connection.Write([]byte("数据格式错误"))
				continue
			}
			r.req = string(contentstr)
			if ProtocolData, err := r.Json2map(); err == nil {
				IMsg.HandleCltProtocol(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData, this.Connection)
			}
		}
	}
}

func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

//func (this *OnlineUser) writeLoop() {
//	defer func() {
//		if err := recover(); err != nil {
//			strerr := fmt.Sprintf("%s", err)
//			glog.Info("异常捕获:", strerr)
//		}
//	}()
//
//	//this.PlayerSendMessage(this.outChan)
//
//	for {
//		select {
//		case data := <-this.outChan:
//			if iret := this.PlayerSendMessage(data); iret == 2 {
//				this.Connection.Close()
//				runtime.Goexit() //new24
//				goto ERR
//			}
//		case <-this.goExit:
//			this.Connection.Close()
//			runtime.Goexit() //new24
//		}
//	}
//ERR:
//	this.Connection.Close()
//	runtime.Goexit()
//}

func (this *OnlineUser) PlayerSendMessage(senddata interface{}) int {

	glog.Info("协程的数量 :", runtime.NumGoroutine())
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err1 := jsoniter.Marshal(senddata)
	if err1 != nil {
		glog.Error("PlayerSendMessage json.Marshal data fail ! err:", err1.Error())
		glog.Flush()
		return 1
	}
	err := websocket.JSON.Send(this.Connection, b)
	if err != nil {
		glog.Error("PlayerSendMessage send data fail ! err:", err.Error())
		glog.Flush()
		return 2
	}
	return 0
}

type Requestbody struct {
	req string
}

func (r *Requestbody) Json2map() (s map[string]interface{}, err error) {
	var result map[string]interface{}
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := jsoniter.Unmarshal([]byte(r.req), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func PlayerSendToServer(conn *websocket.Conn, data interface{}) {
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	jsons, err := jsoniter.Marshal(data)
	if err != nil {
		glog.Info("err:", err.Error())
		return
	}
	errq := websocket.Message.Send(conn, jsons)
	if errq != nil {
		glog.Info(errq)
	}
	return
}

//------------------------------------------------------------------------------
func PlayerSendToProxyServer(conn *websocket.Conn, senddata interface{}, strOpenID string) {
	if len(strOpenID) > 50 {
		return
	}
	data := Proto_Proxy.G2Proxy_SendData{
		Protocol:  1,
		Protocol2: Proto_Proxy.G2Proxy_SendDataProto,
		OpenID:    strOpenID,
		Data:      senddata,
	}
	var sssend interface{}
	sssend = data
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	jsons, err := jsoniter.Marshal(sssend)
	if err != nil {
		glog.Info("err:", err.Error())
		return
	}
	errq := websocket.Message.Send(conn, jsons)
	if errq != nil {
		glog.Info(errq)
	}
	return
}

func PlayerSendMessageOfProxy(conn *websocket.Conn, senddata interface{}, strServerID string) int {

	datasend := Proto_Proxy.C2Proxy_SendData{
		Protocol:  1,
		Protocol2: 1,
		ServerID:  strServerID,
		Data:      senddata,
	}
	var sssend interface{}
	sssend = datasend
	glog.Info("协程的数量 :", runtime.NumGoroutine())
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err1 := jsoniter.Marshal(sssend)
	if err1 != nil {
		glog.Error("PlayerSendMessage json.Marshal data fail ! err:", err1.Error())
		glog.Flush()
		return 1
	}
	glog.Flush()
	encoding := base64.StdEncoding.EncodeToString(b)
	err := websocket.JSON.Send(conn, encoding)
	if err != nil {
		glog.Error("PlayerSendMessage send data fail ! err:", err.Error())
		glog.Flush()
		return 2
	}
	return 0
}

func PlayerSendMessageOfExit(conn *websocket.Conn, senddata interface{}, strServerID string) int {

	var sssend interface{}
	sssend = senddata
	glog.Info("协程的数量 :", runtime.NumGoroutine())
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err1 := jsoniter.Marshal(sssend)
	if err1 != nil {
		glog.Error("PlayerSendMessage json.Marshal data fail ! err:", err1.Error())
		glog.Flush()
		return 1
	}
	glog.Flush()
	encoding := base64.StdEncoding.EncodeToString(b)
	err := websocket.JSON.Send(conn, encoding)
	if err != nil {
		glog.Error("PlayerSendMessage send data fail ! err:", err.Error())
		glog.Flush()
		return 2
	}
	return 0
}

// WebSocketStart websocket启动
func WebSocketStart(url string,
	route string,
	BuildConnection func(ws *websocket.Conn),
	conntype int,
	serverId int,
	proxyUrl []string, //[0] = ProxyHost;[1]=ProxyPort,[2]=ProxyPath
	GameServerReceive func(ws *websocket.Conn),
	ConnXZ *websocket.Conn) {
	var StartDesc = ""
	if conntype == ConnProxy { //作为内网的服务器连接代理服务器
		proxyURL := AddParamsToGetReq("ws", proxyUrl, map[string]string{"data": "{ID:1}"})
		glog.Infof("connect to proxy addr:%s\n", proxyURL)
		conn, err := websocket.Dial(proxyURL, "", "test://golang/")
		if err != nil {
			glog.Errorln("err:", err.Error())
			return
		}
		ConnXZ = conn
		data := Proto_Proxy.G2Proxy_ConnData{
			Protocol:  1,
			Protocol2: Proto_Proxy.G2Proxy_ConnDataProto,
			ServerID:  util.MD5_LollipopGO(strconv.Itoa(serverId)),
		}
		PlayerSendToServer(conn, data)
		go GameServerReceive(conn)
	} else if conntype == StartProxy {
		StartDesc = "proxy server"
	}
	http.Handle("/"+route, websocket.Handler(BuildConnection))
	glog.Infof("game listen to:[%s]\n", route)
	glog.Info("game start ok ", StartDesc)
	if err := http.ListenAndServe(url, nil); err != nil {
		glog.Info("Entry nil", err.Error())
		glog.Flush()
		return
	}
}

//添加参数到get请求
func AddParamsToGetReq(tpType string, strArr []string, paramsMap map[string]string) string {
	urlPath := getUrlPath(tpType, strArr)
	if len(paramsMap) <= 0 || paramsMap == nil { //如果没有参数需要添加直接返回当前路径
		return urlPath
	}
	urlPath = urlPath + "?" //如果参数个数大于等于0,路径后缀加上?
	paramList := make([]string, 0)
	for k, v := range paramsMap {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	tempStr := strings.Join(paramList, "&")
	return fmt.Sprintf("%s%s", urlPath, tempStr)
}

//获取url路径
func getUrlPath(tpType string, strArr []string) string {
	urlPath := strings.Join(strArr, "")
	return fmt.Sprintf("%s://%s", tpType, urlPath)
}
