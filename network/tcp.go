package impl

import (
	MsgHandleClt "LollipopGo/global_Interface"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"strings"
)

// TCP 格式
type OnlineTCP struct {
	Listener   *net.Listener
	Connection *net.Conn
	inChan     chan string
	outChan    chan interface{}
	closeChan  chan int
	goExit     chan int
	isClosed   bool
	HandleClt  MsgHandleClt.Msg_data
}

func Bin()  {
	ln, err := net.Listen("tcp", ":10010")

	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		_=conn
	}
}

// 初始化网络
func InitConnectionTCP(tcpConn *net.Conn,Listener *net.Listener) (*OnlineTCP, error) {

	conn := &OnlineTCP{
		Listener:Listener,
		Connection: tcpConn,
		inChan:     make(chan string, BytebufLen),
	}

	go conn.handleLoop()
	conn.readLoop()

 	return conn, nil
}

func (this *OnlineTCP) readLoop() {

	for {
		go func(conn *net.Conn) {
			var buffer = make([]byte, 1024, 1024)
			for {
				n, e := (*conn).Read(buffer)
				if e != nil {
					if e == io.EOF {
						IMsg.CloseEOF(this.Connection)
						break
					}
					break
				}
				fmt.Println("receive from client:", buffer[:n])
			}
			select {
			case this.inChan <- string(buffer):
			}
		}(this.Connection)
	}
}

func (this *OnlineTCP) handleLoop() {

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
