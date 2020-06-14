package deepcopy

import (
	"LollipopGo/tools/jsonutils"
	"encoding/json"
)

//用json序列化的方法深拷贝，比反射更慢；有时候需要用json的tag去除某些字段充当RO
//这时候可以用该函数
func CopyJsonObject(obj interface{}) jsonutils.JsonObject {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	var cp jsonutils.JsonObject
	err = json.Unmarshal(bytes, &cp)
	return cp
}
