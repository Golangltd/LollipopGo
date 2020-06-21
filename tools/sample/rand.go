package sample

import (
	"LollipopGo/tools/collection"
	"LollipopGo/tools/tz"
	"math/rand"
)

//初始化random
func InitRand(){
	rand.Seed(tz.GetNowTsMs())
}

//随机字符串，包含大小写字母和数字
func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		x := rand.Intn(3)
		switch x {
		case 0:
			bytes[i] = byte(RandInt(65, 90)) //大写字母
		case 1:
			bytes[i] = byte(RandInt(97, 122))
		case 2:
			bytes[i] = byte(rand.Intn(10))
		}
	}
	return string(bytes)
}

//闭区间
func RandInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func RandInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func Shuffle(array []int) {
	for i := range array {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func ShuffleInt32(array []int32) {
	for i := range array {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func ShuffleInt64(array []int64) {
	for i := range array {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func ShuffleUint64(array []uint64) {
	for i := range array {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func RandChoiceInt32(array []int32, n int) []int32 {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []int32{array[rand.Intn(len(array))]}
	}
	tmp := make([]int32, len(array))
	copy(tmp, array)
	if len(tmp) <= n {
		return tmp
	}
	ShuffleInt32(tmp)
	return tmp[:n]
}

func RandChoice(array []int, n int) []int {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []int{array[rand.Intn(len(array))]}
	}
	tmp := make([]int, len(array))
	copy(tmp, array)
	if len(tmp) <= n {
		return tmp
	}
	Shuffle(tmp)
	return tmp[:n]
}

//根据权重随机，返回对应选项的索引，O(n)
func WeightedChoice(weightArray []int) int {
	if weightArray == nil {
		return -1
	}
	total := collection.SumInt(weightArray)
	rv := rand.Int63n(total)
	for i, v := range weightArray {
		if rv < int64(v) {
			return i
		}
		rv -= int64(v)
	}
	return len(weightArray) - 1
}

//是否命中百分比
func HitRate100(rate int) bool {
	return rand.Intn(100) < rate
}

//是否命中千分比
func HitRate1000(rate int) bool {
	return rand.Intn(1000) < rate
}

//是否命中万分比
func HitRate10000(rate int) bool {
	return rand.Intn(10000) < rate
}
