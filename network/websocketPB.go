package impl

import (
	Proto_Proxy "LollipopGo/Proxy_Server/Proto"
	MsgHandleClt "LollipopGo/global_Interface"
	"LollipopGo/util"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"strconv"
)

var BytebufLenPB int = 10000
var IMsgPB MsgHandleClt.Msg_dataPB

type OnlineUserPB struct {
	Connection *websocket.Conn
	inChan     chan []byte
	outChan    chan interface{}
	closeChan  chan int
	goExit     chan int
	isClosed   bool
	HandleClt  MsgHandleClt.Msg_dataPB
}

func InitConnectionPB(wsConn *websocket.Conn) (*OnlineUserPB, error) {
	conn := &OnlineUserPB{
		Connection: wsConn,
		inChan:     make(chan []byte, BytebufLenPB),
	}
	go conn.handleLoopPB()
	conn.readLoopPB()

	return conn, nil
}

func (this *OnlineUserPB) readLoopPB() {

	for {
		var content []byte
		err := websocket.Message.Receive(this.Connection, &content)
		if err != nil {
			if err == io.EOF || err == io.ErrClosedPipe || len(content) == 0 || err == io.ErrNoProgress {
				IMsgPB.CloseEOFPB(this.Connection)
				return
			}
			break
		}
		select {
		case this.inChan <- content:
		}
	}
}

func (this *OnlineUserPB) handleLoopPB() {

	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			glog.Info("异常捕获:", strerr)
		}
	}()

	for {
		var r RequestbodyPB
		select {
		case r.req = <-this.inChan:
		}
		if len(r.req) <= 0 {
			continue
		}
		if ProtocolData, err := r.Json2mapPB(); err == nil {
			IMsgPB.HandleCltProtocolPB(Proto_Proxy.Proxy_CMD(ProtocolData.Protocol), Proto_Proxy.Proxy_CMD(ProtocolData.Protocol2), ProtocolData.PackageData, this.Connection)
		} else if ProtocolDataServer, err := r.Json2mapPBServer(); err == nil {
			IMsgPB.HandleCltProtocolPB(Proto_Proxy.Proxy_CMD(ProtocolDataServer.Protocol), Proto_Proxy.Proxy_CMD(ProtocolDataServer.Protocol2), ProtocolDataServer.PackageData, this.Connection)
		}
	}
}

type RequestbodyPB struct {
	req []byte
}

func (r *RequestbodyPB) Json2mapPB() (s *Proto_Proxy.ProxyC2S_SendData, err error) {
	newClientData := &Proto_Proxy.ProxyC2S_SendData{}
	err = proto.Unmarshal(r.req, newClientData)
	if err != nil {
		return &Proto_Proxy.ProxyC2S_SendData{}, err
	}
	return newClientData, nil
}

func (r *RequestbodyPB) Json2mapPBServer() (s *Proto_Proxy.ProxyS2C_SendData, err error) {
	newClientData := &Proto_Proxy.ProxyS2C_SendData{}
	err = proto.Unmarshal(r.req, newClientData)
	if err != nil {
		return &Proto_Proxy.ProxyS2C_SendData{}, err
	}
	return newClientData, nil
}

// WebSocketStart websocket启动
func WebSocketStartPB(url string,
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

func PlayerSendToServerPB(conn *websocket.Conn, data []byte) {
	errq := websocket.Message.Send(conn, data)
	if errq != nil {
		glog.Info(errq)
	}
	return
}

//------------------------------------------------------------------------------
func PlayerSendToProxyServerPB(conn *websocket.Conn, senddata []byte, strOpenID string) {
	if len(strOpenID) > 50 {
		return
	}
	data := &Proto_Proxy.ProxyS2C_SendData{
		Protocol:    1,
		Protocol2:   int32(Proto_Proxy.Proxy_P2C_SendData),
		OpenId:      strOpenID,
		PackageData: senddata,
	}

	PackageData, err := proto.Marshal(data)
	if err != nil {
		glog.Info("序列化失败:", err)
		return
	}

	errq := websocket.Message.Send(conn, PackageData)
	if errq != nil {
		glog.Info(errq)
	}
	return
}
