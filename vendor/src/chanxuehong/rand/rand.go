package rand

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"
)

const (
	saltLen            = 45   // see New(), 6+4+45==55<56, Best performance for md5
	saltUpdateInterval = 3600 // seconds
)

var (
	gSalt     = make([]byte, saltLen)
	gSequence = Uint32()

	gMutex                   sync.Mutex
	gSaltLastUpdateTimestamp int64 = -saltUpdateInterval
)

// New returns 16-byte raw random bytes.
// It is not printable, you can use encoding/hex or encoding/base64 to print it.
func New() (rd [16]byte) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()

	if timeNowUnix >= atomic.LoadInt64(&gSaltLastUpdateTimestamp)+saltUpdateInterval {
		gMutex.Lock() // Lock
		if timeNowUnix >= gSaltLastUpdateTimestamp+saltUpdateInterval {
			gSaltLastUpdateTimestamp = timeNowUnix
			gMutex.Unlock() // Unlock

			Read(gSalt)
			copy(rd[:], gSalt)
			return
		}
		gMutex.Unlock() // Unlock
	}
	sequence := atomic.AddUint32(&gSequence, 1)

	timeNowUnixNano := timeNow.UnixNano()
	var src [6 + 4 + saltLen]byte // 6+4+45==55
	src[0] = byte(timeNowUnixNano >> 40)
	src[1] = byte(timeNowUnixNano >> 32)
	src[2] = byte(timeNowUnixNano >> 24)
	src[3] = byte(timeNowUnixNano >> 16)
	src[4] = byte(timeNowUnixNano >> 8)
	src[5] = byte(timeNowUnixNano)
	src[6] = byte(sequence >> 24)
	src[7] = byte(sequence >> 16)
	src[8] = byte(sequence >> 8)
	src[9] = byte(sequence)
	copy(src[10:], gSalt)

	return md5.Sum(src[:])
}

// NewHex returns 32-byte hex-encoded bytes.
func NewHex() (rd []byte) {
	rdx := New()
	rd = make([]byte, hex.EncodedLen(len(rdx)))
	hex.Encode(rd, rdx[:])
	return
}
