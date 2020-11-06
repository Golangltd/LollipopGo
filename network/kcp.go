package impl

import (
	MsgHandleClt "LollipopGo/global_Interface"
	"fmt"
	"github.com/golang/glog"
	"github.com/xtaci/kcp-go"
	"io"
	"net"
	"strings"
)

// KCP 格式
type OnlineKCP struct {
	Listener   *kcp.Listener
	Connection *kcp.UDPSession
	inChan     chan string
	outChan    chan interface{}
	closeChan  chan int
	goExit     chan int
	isClosed   bool
	HandleClt  MsgHandleClt.Msg_data
}
// 初始化网络
func InitConnectionKCP(kcpConn *kcp.UDPSession,Listener*kcp.Listener) (*OnlineKCP, error) {

	conn := &OnlineKCP{
		Listener:Listener,
		Connection: kcpConn,
		inChan:     make(chan string, BytebufLen),
	}

	go conn.handleLoop()
	conn.readLoop()

	return conn, nil
}

func (this *OnlineKCP) readLoop() {

	for {
		go func(conn net.Conn) {
			var buffer = make([]byte, 1024, 1024)
			for {
				n, e := conn.Read(buffer)
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

func (this *OnlineKCP) handleLoop() {

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
