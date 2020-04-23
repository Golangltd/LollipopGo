package network

import (
	"FenDZ/glog-master"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	lenMsgLen    int
	minMsgLen    uint32
	maxMsgLen    uint32
	littleEndian bool
}

// 默认在没有设置客户端信息的时候
// 消息的长度
// 每一个客户端的参数
func NewMsgParser() *MsgParser {
	p := new(MsgParser)
	p.lenMsgLen = 2
	p.minMsgLen = 1
	p.maxMsgLen = 4096
	p.littleEndian = false

	return p
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}

	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// --------------
// | len | data |
// --------------
//type MsgParser struct {
//	lenMsgLen    int
//	minMsgLen    uint32
//	maxMsgLen    uint32
//	littleEndian bool
//}

// goroutine safe
func (p *MsgParser) Read(conn *TCPConn) ([]byte, error) {

	glog.Info("服务器接受到的数据结构<1>:", p)
	glog.Info("服务器接受到的数据结构<2>:", p.lenMsgLen)
	glog.Info("服务器接受到的数据结构<3>:", p.littleEndian)

	var b [4]byte
	bufMsgLen := b[:p.lenMsgLen]
	glog.Info("服务器接受到的数据结构<n>bufMsgLen:", bufMsgLen)

	// read len
	// 客户端退出后 打印
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}

	// parse len
	var msgLen uint32
	switch p.lenMsgLen {
	case 1:
		msgLen = uint32(bufMsgLen[0])
	case 2:
		if p.littleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
			glog.Info("服务器接受到的数据结构<4>:", msgLen)
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
			glog.Info("服务器接受到的数据结构<4>:", msgLen)
		}
	case 4:
		if p.littleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}

	// check len
	glog.Info("服务器接受到的数据结构<3>-----:", msgLen)
	glog.Info("服务器接受到的数据结构<3>------:", p.maxMsgLen)
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}

	// data
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		glog.Info("服务器接受到的数据结构<3>------:", err)
		return nil, err
	}
	glog.Info("服务器接受到的数据结构<3>------:", msgData)

	return msgData, nil
}

// goroutine safe
func (p *MsgParser) Write(conn *TCPConn, args ...[]byte) error {
	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	msg := make([]byte, uint32(p.lenMsgLen)+msgLen)

	// write len
	switch p.lenMsgLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, msgLen)
		} else {
			binary.BigEndian.PutUint32(msg, msgLen)
		}
	}

	// write data
	l := p.lenMsgLen
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	conn.Write(msg)

	return nil
}
