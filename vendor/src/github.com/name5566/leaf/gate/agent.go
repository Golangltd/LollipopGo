package gate

import (
	"FenDZ/go-concurrentMap-master"
	"net"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}

// 在线玩家的数据的结构体
type OnlineUser struct {
	Connection Agent                     // 链接的信息
	StrMD5     string                    // 用的UID标示
	MapSafe    *concurrent.ConcurrentMap // 并发安全的map
}
