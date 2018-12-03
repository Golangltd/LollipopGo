package util

import (
	"crypto/md5"
	"encoding/hex"
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
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//------------------------------------------------------------------------------
