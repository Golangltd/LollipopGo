package DFA

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	FILE_FILTER = "filter.txt"
)

var (
	ConfExample *ConfigFilter
)

type ConfigFilter struct {
	FilterList map[rune]*FilterModel //屏蔽字树
}

//加载词库
func InitConfigFilter(configpath string) *ConfigFilter {
	result := new(ConfigFilter)
	{
		li := make(map[rune]*FilterModel)
		//我这里用的是一个文本文件，一行表示一个屏蔽词
		file, err := os.Open(path.Join(configpath, FILE_FILTER))
		if err != nil {
			panic(err)
		}
		barr, _ := ioutil.ReadAll(file)
		bstr := string(barr)
		bstr = strings.ReplaceAll(bstr, "\r", "")
		rows := strings.Split(bstr, "\n")
		for _, row := range rows {
			rowr := []rune(row)
			fmd, ok := li[rowr[0]]
			if !ok {
				fmd = new(FilterModel)
				fmd.NodeStr = rowr[0]
				fmd.Subli = make(map[rune]*FilterModel)
				li[rowr[0]] = fmd
			}
			fmd.IsEnd = filterFor(fmd.Subli, rowr, 1)
		}
		result.FilterList = li
	}
	return result
}

func filterFor(li map[rune]*FilterModel, rowr []rune, index int) bool {
	if len(rowr) <= index {
		return true
	}
	fmd, ok := li[rowr[index]]
	if !ok {
		fmd = new(FilterModel)
		fmd.NodeStr = rowr[index]
		fmd.Subli = make(map[rune]*FilterModel)
		li[rowr[index]] = fmd
	}
	index++
	fmd.IsEnd = filterFor(fmd.Subli, rowr, index)
	return false
}

//屏蔽字结构
type FilterModel struct {
	NodeStr rune //内容
	Subli map[rune]*FilterModel //屏蔽子集合
	IsEnd bool //是否为结束
}

//屏蔽字操作，这个方法就是外部调用的入口方法
func LollipopGoFilterCheck(data string) (result string) {
	filterli := ConfExample.FilterList
	arr := []rune(data)
	for i := 0; i < len(arr); i++ {
		fmd, ok := filterli[arr[i]]
		if !ok {
			continue
		}
		if ok, index := filterChackFor(arr, i+1, fmd.Subli); ok {
			arr[i] = rune('*')
			i = index
		}
	}
	return string(arr)
}

//递归调用检查屏蔽字
func filterChackFor(arr []rune, index int, filterli map[rune]*FilterModel) (bool, int) {
	if len(arr) <= index {
		return false, index
	}
	if arr[index] == rune(' ') {
		if ok, i := filterChackFor(arr, index+1, filterli); ok {
			arr[index] = rune('*')
			return true, i
		}
	}
	fmd, ok := filterli[arr[index]]
	if !ok {
		return false, index
	}
	if fmd.IsEnd {
		arr[index] = rune('*')
		ok, i := filterChackFor(arr, index+1, fmd.Subli)
		if ok {
			return true, i
		}
		return true, index
	} else if ok, i := filterChackFor(arr, index+1, fmd.Subli); ok {
		arr[index] = rune('*')
		return true, i
	}
	return false, index
}
