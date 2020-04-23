package rand

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"time"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// Read reads len(p)-byte raw random bytes to p.
func Read(p []byte) {
	if len(p) <= 0 {
		return
	}

	// get from crypto/rand
	if _, err := cryptorand.Read(p); err == nil {
		return
	}

	// get from math/rand
	timeNowNano := int64(time.Now().Nanosecond())
	timeNowNano = timeNowNano<<32 | timeNowNano
	for len(p) > 0 {
		n := mathrand.Int63() ^ timeNowNano
		switch len(p) {
		case 8:
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			p[5] = byte(n >> 16)
			p[6] = byte(n >> 8)
			p[7] = byte(n)
			return
		case 4:
			p[0] = byte(n >> 24)
			p[1] = byte(n >> 16)
			p[2] = byte(n >> 8)
			p[3] = byte(n)
			return
		case 1:
			p[0] = byte(n)
			return
		case 2:
			p[0] = byte(n >> 8)
			p[1] = byte(n)
			return
		case 3:
			p[0] = byte(n >> 16)
			p[1] = byte(n >> 8)
			p[2] = byte(n)
			return
		case 5:
			p[0] = byte(n >> 32)
			p[1] = byte(n >> 24)
			p[2] = byte(n >> 16)
			p[3] = byte(n >> 8)
			p[4] = byte(n)
			return
		case 6:
			p[0] = byte(n >> 40)
			p[1] = byte(n >> 32)
			p[2] = byte(n >> 24)
			p[3] = byte(n >> 16)
			p[4] = byte(n >> 8)
			p[5] = byte(n)
			return
		case 7:
			p[0] = byte(n >> 48)
			p[1] = byte(n >> 40)
			p[2] = byte(n >> 32)
			p[3] = byte(n >> 24)
			p[4] = byte(n >> 16)
			p[5] = byte(n >> 8)
			p[6] = byte(n)
			return
		default: // len(p) > 8
			_ = p[8]
			p[0] = byte(n >> 56)
			p[1] = byte(n >> 48)
			p[2] = byte(n >> 40)
			p[3] = byte(n >> 32)
			p[4] = byte(n >> 24)
			p[5] = byte(n >> 16)
			p[6] = byte(n >> 8)
			p[7] = byte(n)
			p = p[8:]
		}
	}
}
