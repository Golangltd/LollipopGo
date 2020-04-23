package utf8

import (
	"code.google.com/p/mahonia"
)

// GBK to UTF-8
func GBKConvertUTF8(content string) (ret string) {
	dec := mahonia.NewDecoder("gbk")
	ret, ok := dec.ConvertStringOK(content)
	if ok {
		return ret
	} else {
		return ""
	}
}

// UTF-8 to GBK
func UTF8ConvertGBK(content string) (ret string) {
	enc := mahonia.NewEncoder("gbk")
	ret, ok := enc.ConvertStringOK(content)
	if ok {
		return ret
	} else {
		return ""
	}
}
