package Mysyl_DB

import (
	"LollipopGo/LollipopGo/util"
)

/*
   数据库的通用函数
*/

func CheckErr(err error) {
	util.CheckErr_LollipopGO(err)
}

func GetTime() string {
	return util.GetTime_LollipopGO()
}

func GetNowtimeMD5() string {
	return util.GetNowtimeMD5_LollipopGO()
}

func GetMD5Hash(text string) string {
	return util.MD5_LollipopGO(text)
}
