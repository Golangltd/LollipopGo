package rand

import (
	"encoding/binary"
)

const (
	uint31Mask = ^(uint32(1) << 31)
	uint63Mask = ^(uint64(1) << 63)
)

// Int31 returns a non-negative random 31-bit integer as an int32.
func Int31() int32 {
	u32 := Uint32()
	return int32(u32 & uint31Mask)
}

// Int63 returns a non-negative random 63-bit integer as an int64.
func Int63() int64 {
	u64 := Uint64()
	return int64(u64 & uint63Mask)
}

// Uint32 returns a random 32-bit integer as a uint32.
func Uint32() uint32 {
	var x [4]byte
	Read(x[:])
	return binary.BigEndian.Uint32(x[:])
}

// Uint64 returns a random 64-bit integer as a uint64.
func Uint64() uint64 {
	var x [8]byte
	Read(x[:])
	return binary.BigEndian.Uint64(x[:])
}
