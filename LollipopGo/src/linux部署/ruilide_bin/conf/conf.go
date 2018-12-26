package conf

import (
	"log"
	"time"
)

var (
	// log conf
	LogFlag = log.LstdFlags

	// gate conf // 网关配置
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096 // 消息的长度
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// skeleton conf  框架配置
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)
