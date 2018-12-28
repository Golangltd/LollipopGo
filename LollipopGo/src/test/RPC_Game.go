package main

import (
	"fmt"
)

// 网络接口
type Conn interface {
	// ConnGateWayServer(data interface{})
	// PlayerSendMessage(data interface{})
	// HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{})
	// HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{})
	// Close()
	Destroy()
}

// 每个类型服务器结构信息 实例1
type CST struct {
	data Conn   //  接口形式  ---- 主要是【注册】形式,所有的都可以处理
	UID  string //  处理
}

// 注册数据
func Register(conn Conn) {

	return
}

// 调用
func init() {
	data := &CST{
		UID: "888",
	}
	Register(data)
	//NewData(data)

	return
}

// 不用的实现形式
func (this *CST) Destroy() {
	fmt.Println("111", this.UID)
	return
}

func main() {
	return
}

// new 数据
func NewData(data Conn) {
	data.Destroy()
	return
}
