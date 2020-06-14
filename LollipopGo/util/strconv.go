package util

import (
	"fmt"
	"strconv"
)

func Str2intLollipopgo(data string) int {
	v, err := strconv.Atoi(data)
	if err != nil {
		return -1
	}
	return v
}

func Int2str_LollipopGo(data int) string {
	return strconv.Itoa(data)
}

func toString(arg interface{}) string {
	switch arg.(type) {
	case bool:
		return boolToString(arg.(bool))
	case float32:
		return floatToString(float64(arg.(float32)))
	case float64:
		return floatToString(arg.(float64))
		//case complex64:
		//  p.fmtComplex(complex128(f), 64, verb)
		//case complex128:
		//  p.fmtComplex(f, 128, verb)
	case int:
		return intToString(int64(arg.(int)))
	case int8:
		return intToString(int64(arg.(int8)))
	case int16:
		return intToString(int64(arg.(int16)))
	case int32:
		return intToString(int64(arg.(int32)))
	case int64:
		return intToString(int64(arg.(int64)))
	default:
		return fmt.Sprint(arg)
	}
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'E', -1, 64)
}
func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func boolToString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
