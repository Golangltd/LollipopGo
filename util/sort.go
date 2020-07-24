package util



/*//------------------------------------------------------------------------------
// 例子已经写在简书：https://www.jianshu.com/p/e30a9db07da0
// 详见《彬哥Go语言笔记》
func Sort_LollipopGo(data map[string]*conf.DSQ_Exp, iExp int) int {

	if iExp == 0 {
		return 0
	}
	var length = len(data)
	var ssort []int

	for _, v := range data {
		ssort = append(ssort, Str2intLollipopgo(v.Exp))
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
	return 0
}
*/

//  i := Minimum(1, 3, 5, 7, 9, 10, -1, 1).(int)
func Minimum(first interface{}, rest ...interface{}) interface{} {
	minimum := first

	for _, v := range rest {
		switch v.(type) {
		case int:
			if v := v.(int); v < minimum.(int) {
				minimum = v
			}
		case float64:
			if v := v.(float64); v < minimum.(float64) {
				minimum = v
			}
		case string:
			if v := v.(string); v < minimum.(string) {
				minimum = v
			}
		}
	}
	return minimum
}
