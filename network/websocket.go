package impl

import (
	"LollipopGo/Proxy_Server/Proto"
	"LollipopGo/global_Interface"
	"encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"github.com/json-iterator/go"
	"golang.org/x/net/websocket"
	"io"
	"runtime"
	"strings"
)

var BytebufLen int = 10000
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
	go conn.handleLoop()
	conn.readLoop()

	return conn, nil
}

func (this *OnlineUser) readLoop() {

	for {
		var content string
		err := websocket.Message.Receive(this.Connection, &content)
		if err != nil {
			if err == io.EOF {
				IMsg.CloseEOF(this.Connection)
				break
			}
			break
		}
		select {
		case this.inChan <- content:
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
		var r Requestbody
		select {
		case r.req = <-this.inChan:
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

func (this *OnlineUser) writeLoop() {
	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			glog.Info("异常捕获:", strerr)
		}
	}()
	for {
		select {
		case data := <-this.outChan:
			if iret := this.PlayerSendMessage(data); iret == 2 {
				this.Connection.Close()
				runtime.Goexit() //new24
				goto ERR
			}
		case <-this.goExit:
			this.Connection.Close()
			runtime.Goexit() //new24
		}
	}
ERR:
	this.Connection.Close()
	runtime.Goexit()
}

func (this *OnlineUser) PlayerSendMessage(senddata interface{}) int {

	glog.Info("协程的数量 :", runtime.NumGoroutine())
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	/*b, err1 := json.Marshal(senddata)*/
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
		//glog.Info("Json2map:", err.Error())
		return nil, err
	}
	/*	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
		//glog.Info("Json2map:", err.Error())
		return nil, err
	}*/
	return result, nil
}

func PlayerSendToServer(conn *websocket.Conn, data interface{}) {
	var jsoniter = jsoniter.ConfigCompatibleWithStandardLibrary
	jsons, err := jsoniter.Marshal(data)
	/*jsons, err := json.Marshal(data)*/
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
	/*jsons, err := json.Marshal(sssend)*/
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
	/*	b, err1 := json.Marshal(sssend)*/
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
	/*	b, err1 := json.Marshal(sssend)*/
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
