package num

import "math/rand"

//是否命中百分比
func HitRate100(rate int32) bool {
	return rand.Int31n(100) < rate
}

//是否命中千分比
func HitRate1000(rate int32) bool {
	return rand.Int31n(1000) < rate
}

//是否命中万分比
func HitRate10000(rate int32) bool {
	return rand.Int31n(10000) < rate
}