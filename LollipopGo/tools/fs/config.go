package fs

import (
	"sync"
)

//解析器需要实现的接口
type IConfigParser interface {
	ReloadConfig(path string, init bool) bool //重载配置
	GetConfig() interface{}                   //获取配置
}

//解析器的默认实现，用于嵌套
type ParserMixIn struct {
	sync.RWMutex
	lastModified map[string]int64
}

//SetLastModifyTime update lastModifyTime
func (pmi *ParserMixIn) SetLastModifyTime(path string, ts int64) {
	if pmi.lastModified == nil {
		pmi.lastModified = make(map[string]int64)
	}
	pmi.lastModified[path] = ts
}

//GetPathLastModifyTime
func (pmi *ParserMixIn) GetPathModifyTime(path string) int64 {
       if pmi.lastModified == nil {
        return 0
       }
    return pmi.lastModified[path]
}

//CheckModify return if modified and last modify time
func (pmi *ParserMixIn) CheckModify(path string) (bool, int64) {
	ts, err := GetLastModifyTime(path)
	if err != nil {
		return false, 0
	}
	if pmi.lastModified == nil {
		return true, ts
	}
	return ts != pmi.lastModified[path], ts
}

//周期性监测文件变化，调用Parser的回调
func WatchConfigFiles(parsers map[string]IConfigParser) {
	for path, parser := range parsers {
		if parser.ReloadConfig(path, false) {
		}
	}
}
