package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//------------------------------------------------------------------------------

// package main

// import (
//     "fmt"
//     "strconv"
//     "time"
// )

// func main() {
//     t := time.Now()
//     fmt.Println(t)

//     fmt.Println(t.UTC().Format(time.UnixDate))

//     fmt.Println(t.Unix())

//     timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
//     fmt.Println(timestamp)
//     timestamp = timestamp[:10]
//     fmt.Println(timestamp)
// }

// 输出：
// 2017-06-21 11:52:29.0826692 + 0800 CST
// Wed Jun 21 03:52:29 UTC 2017
// 1498017149
// 1498017149082669200
// 1498017149

// 生成时间戳的函数
func UTCTime_LollipopGO() string {
	t := time.Now()
	return strconv.FormatInt(t.UTC().UnixNano(), 10)
}

//------------------------------------------------------------------------------
// package main

// import (
//     "crypto/md5"
//     "encoding/hex"
//     "fmt"
// )

// func main() {
//     h := md5.New()
//     h.Write([]byte("123456")) // 需要加密的字符串为 123456
//     cipherStr := h.Sum(nil)
//     fmt.Println(cipherStr)
//     fmt.Printf("%s\n", hex.EncodeToString(cipherStr)) // 输出加密结果
// }

// MD5 实现 :主要是针对 字符串的加密
func MD5_LollipopGO(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//------------------------------------------------------------------------------
//返回[0,max)的随机整数
func Randnum_LollipopGO(max int) int {

	if max == 0 {
		panic("随机函数，传递参数错误!")
		return -1
	}
	// 随机种子:系统时间
	rand.Seed(time.Now().Unix())
	return rand.Intn(max)
}

//------------------------------------------------------------------------------

func CheckErr_LollipopGO(err error) {
	if err != nil {
		panic(err)
		fmt.Println("err:", err)
	}
}

func GetTime_LollipopGO() string {
	const shortForm = "2006-01-02 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}

func GetNowtimeMD5_LollipopGO() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return MD5_LollipopGO(timestamp)
}
