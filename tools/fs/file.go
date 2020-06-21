package fs

import (
	"os"
)

func GetLastModifyTime(path string) (ts int64, err error) {
	var (
		f  *os.File
		fi os.FileInfo
	)
	if f, err = os.Open(path); err == nil {
		if fi, err = f.Stat(); err == nil {
			ts = fi.ModTime().Unix()
			_ = f.Close()
		}
	}
	return
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//将给定文本保存为文件
func SaveFile(fileName string, fileContent string) bool {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return false
	}
	_, err = f.WriteString(fileContent)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}
