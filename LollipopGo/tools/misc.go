package tools

import (
	"LollipopGo/tools/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	configDir = "_config"
)

//约定程序二进制文件和_config文件夹在同一层
//否则逐层往上找直到找到（方便进行单元测试）
func GetConfigDir() string {
	wd, _ := os.Getwd()

	x := filepath.Join(wd, configDir)
	for !fs.Exists(x) {
		if wd == "/" {
		}
		wd = filepath.Dir(wd)
		x = filepath.Join(wd, configDir)
	}
	return x
}

type TextDuration struct {
	time.Duration
}

func (d *TextDuration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

//通用recover函数，在单独协程的最开始使用defer调用
func RecoverFromPanic(cb func()) {
	if r := recover(); r != nil {
		if cb != nil {
			cb()
		}
	}
}
