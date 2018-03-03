/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月2日
*/

package cache

import (
	"errors"
)

var (
	ErrKeyNotFound           = errors.New("Key not found in cache")
	ErrKeyNotFoundOrLoadable = errors.New("Key not found and could not be loaded into cache")
)
