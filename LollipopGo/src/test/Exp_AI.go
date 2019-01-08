package main

import (
	"LollipopGo/LollipopGo/util"
	"fmt"
)

var G_Exp_Lev map[string]*EXP

// 经验表结构
type EXP struct {
	Lev string
	Exp string
}

func init() {
	G_Exp_Lev = make(map[string]*EXP)
}

func SaveDataslice(data map[string]*EXP, iExp int) int {
	var length = len(data)
	var ssort []int

	for _, v := range data {
		ssort = append(ssort, util.Str2int_LollipopGo(v.Exp))
	}

	for i := 1; i < length; i++ {
		for j := i; j > 0 && ssort[j] < ssort[j-1]; j-- {
			ssort[j], ssort[j-1] = ssort[j-1], ssort[j]
		}
	}
	fmt.Println(ssort)
	for index, val := range ssort {
		//fmt.Printf("index array[%d] = %d\n", index, val)
		if iExp == val {
			return index
		}
	}
	return 0
}

func main() {

	for i := 1; i <= 10; i++ {
		// 模拟经验表配置数据
		data := &EXP{
			Lev: util.Int2str_LollipopGo(i),
			Exp: util.Int2str_LollipopGo(i * 10),
		}
		G_Exp_Lev[util.Int2str_LollipopGo(i)] = data
	}

	data := &EXP{
		Lev: util.Int2str_LollipopGo(11111),
		Exp: util.Int2str_LollipopGo(55),
	}
	G_Exp_Lev[util.Int2str_LollipopGo(11)] = data
	sss := SaveDataslice(G_Exp_Lev, 55)
	fmt.Println("排名：", sss)
}
