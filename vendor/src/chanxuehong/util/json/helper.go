package json

import "unsafe"

func ConvertToStdIntSlice(s []Int) []int {
	return *((*[]int)(unsafe.Pointer(&s)))
}

func ConvertToIntSlice(s []int) []Int {
	return *((*[]Int)(unsafe.Pointer(&s)))
}

func ConvertToStdUintSlice(s []Uint) []uint {
	return *((*[]uint)(unsafe.Pointer(&s)))
}

func ConvertToUintSlice(s []uint) []Uint {
	return *((*[]Uint)(unsafe.Pointer(&s)))
}

func ConvertToStdInt64Slice(s []Int64) []int64 {
	return *((*[]int64)(unsafe.Pointer(&s)))
}

func ConvertToInt64Slice(s []int64) []Int64 {
	return *((*[]Int64)(unsafe.Pointer(&s)))
}

func ConvertToStdUint64Slice(s []Uint64) []uint64 {
	return *((*[]uint64)(unsafe.Pointer(&s)))
}

func ConvertToUint64Slice(s []uint64) []Uint64 {
	return *((*[]Uint64)(unsafe.Pointer(&s)))
}

func ConvertToStdFloat64Slice(s []Float64) []float64 {
	return *((*[]float64)(unsafe.Pointer(&s)))
}

func ConvertToFloat64Slice(s []float64) []Float64 {
	return *((*[]Float64)(unsafe.Pointer(&s)))
}
