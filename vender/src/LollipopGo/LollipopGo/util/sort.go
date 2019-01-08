package util

import (
	"LollipopGo/LollipopGo/conf"
)

//------------------------------------------------------------------------------
// 例子已经写在简书：https://www.jianshu.com/p/e30a9db07da0
// 详见《彬哥Go语言笔记》
func Sort_LollipopGo(data map[string]*conf.DSQ_Exp, iExp int) int {
	var length = len(data)
	var ssort []int

	for _, v := range data {
		ssort = append(ssort, Str2int_LollipopGo(v.Exp))
	}

	for i := 1; i < length; i++ {
		for j := i; j > 0 && ssort[j] < ssort[j-1]; j-- {
			ssort[j], ssort[j-1] = ssort[j-1], ssort[j]
		}
	}
	for index, val := range ssort {
		if iExp == val {
			return index
		}
	}
	panic("排序出错")
	return 0
}
