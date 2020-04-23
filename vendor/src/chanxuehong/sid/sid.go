package sid

import (
	"crypto/sha1"
	"encoding/base64"
	"os"
	"sync"
	"time"

	"github.com/chanxuehong/internal"
	"github.com/chanxuehong/rand"
)

//   56bits unix100ns + 12bits pid + 12bits sequence + 48bits node + 64bits hashsum
//
//   +------ 0 ------+------ 1 ------+------ 2 ------+------ 3 ------+
//   +0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1+
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                          time_low                             |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |       time_mid                |        time_hi_and_pid_low    |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |clk_seq_hi_pid |  clk_seq_low  |         node (0-1)            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                         node (2-5)                            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                         hash (0-3)                            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                         hash (4-7)                            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

var pid = hash(uint64(os.Getpid())) // 12-bit hashsum of os.Getpid(), read only.

// hash uint64 to a 12-bit integer value.
func hash(x uint64) uint64 {
	return (x ^ x>>12 ^ x>>24 ^ x>>36 ^ x>>48 ^ x>>60) & 0xfff
}

var node = internal.MAC[:] // read only

const (
	sequenceMask       = 0xfff // 12bits
	saltLen            = 43    // see New(), 8+4+43==55<56, Best performance for sha1.
	saltUpdateInterval = 3600  // seconds
)

var (
	gSalt = make([]byte, saltLen)

	gMutex sync.Mutex // protect following

	gSequenceStart uint32 = rand.Uint32() & sequenceMask
	gLastTimestamp int64  = -1
	gLastSequence  uint32 = gSequenceStart

	gSaltLastUpdateTimestamp int64  = -saltUpdateInterval
	gSaltSequence            uint32 = rand.Uint32()
)

// New returns a unique 32-byte url-safe string.
func New() string {
	var (
		timeNow     = time.Now()
		timeNowUnix = timeNow.Unix()

		timestamp = unix100nano(timeNow)
		sequence  uint32

		saltShouldUpdate = false
		saltSequence     uint32
	)

	gMutex.Lock() // Lock
	switch {
	case timestamp > gLastTimestamp:
		sequence = gSequenceStart
		gLastTimestamp = timestamp
		gLastSequence = sequence
	case timestamp == gLastTimestamp:
		sequence = (gLastSequence + 1) & sequenceMask
		if sequence == gSequenceStart {
			timestamp = tillNext100nano(timestamp)
			gLastTimestamp = timestamp
		}
		gLastSequence = sequence
	default:
		gSequenceStart = rand.Uint32() & sequenceMask // NOTE
		sequence = gSequenceStart
		gLastTimestamp = timestamp
		gLastSequence = sequence
	}
	if timeNowUnix >= gSaltLastUpdateTimestamp+saltUpdateInterval {
		saltShouldUpdate = true
		gSaltLastUpdateTimestamp = timeNowUnix
	}
	gSaltSequence++
	saltSequence = gSaltSequence
	gMutex.Unlock() // Unlock

	// 56bits unix100ns + 12bits pid + 12bits sequence + 48bits node + 64bits hashsum
	var idx [24]byte

	// time_low
	idx[0] = byte(timestamp >> 24)
	idx[1] = byte(timestamp >> 16)
	idx[2] = byte(timestamp >> 8)
	idx[3] = byte(timestamp)

	// time_mid
	idx[4] = byte(timestamp >> 40)
	idx[5] = byte(timestamp >> 32)

	// time_hi_and_pid_low
	idx[6] = byte(timestamp >> 48)
	idx[7] = byte(pid)

	// clk_seq_hi_pid
	idx[8] = byte(sequence>>8) & 0x0f
	idx[8] |= byte(pid>>8) << 4

	// clk_seq_low
	idx[9] = byte(sequence)

	// node
	copy(idx[10:], node)

	// hashsum
	if saltShouldUpdate {
		rand.Read(gSalt)
		copy(idx[16:], gSalt)
	} else {
		var src [8 + 4 + saltLen]byte // 8+4+43==55

		src[0] = byte(timestamp >> 56)
		src[1] = byte(timestamp >> 48)
		src[2] = byte(timestamp >> 40)
		src[3] = byte(timestamp >> 32)
		src[4] = byte(timestamp >> 24)
		src[5] = byte(timestamp >> 16)
		src[6] = byte(timestamp >> 8)
		src[7] = byte(timestamp)
		src[8] = byte(saltSequence >> 24)
		src[9] = byte(saltSequence >> 16)
		src[10] = byte(saltSequence >> 8)
		src[11] = byte(saltSequence)
		copy(src[12:], gSalt)

		hashsum := sha1.Sum(src[:])
		copy(idx[16:], hashsum[:])
	}

	id := make([]byte, 32)
	base64.URLEncoding.Encode(id, idx[:])
	return string(id)
}
