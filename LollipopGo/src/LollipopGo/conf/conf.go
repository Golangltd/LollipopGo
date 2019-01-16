package conf

import (
	"log"
	"time"
)

var (
	LogFlag                   = log.LstdFlags
	PendingWriteNum           = 2000
	MaxMsgLen          uint32 = 4096
	HTTPTimeout               = 10 * time.Second
	LenMsgLen                 = 2
	LittleEndian              = false
	GoLen                     = 10000
	TimerDispatcherLen        = 10000
	AsynCallLen               = 10000
	ChanRPCLen                = 10000
)
