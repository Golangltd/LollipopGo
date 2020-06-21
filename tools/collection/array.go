package collection

func SumInt32s(array []int32) int64 {
	var sum int64
	for _, v := range array {
		sum += int64(v)
	}
	return sum
}

func SumInt(array []int) int64 {
	var sum int64
	for _, v := range array {
		sum += int64(v)
	}
	return sum
}

//会删去所有相等的元素
func DeleteInt32s(array []int32, elem ...int32) []int32 {
	toDelete := NewInt32Set(elem...)
	result := make([]int32, 0, len(array))
	for _, v := range array {
		if !toDelete.Contains(v) {
			result = append(result, v)
		}
	}
	return result
}

//只会删除第一个相等的元素
func DeleteInt32(array []int32, elem int32) []int32 {
	index := -1
	for i, v := range array {
		if v == elem {
			index = i
			break
		}
	}
	if index == -1 {
		return array
	}
	return DeleteInt32ByIndex(array, index)
}

func DeleteInt32ByIndex(array []int32, index int) []int32 {
	return append(array[:index], array[index+1:]...)
}

func GetElementIndexInt32(array []int32, elem int32) int {
	for i, d := range array {
		if d == elem {
			return i
		}
	}
	return -1
}

//是否完全包含elem，如果elem中有重复的元素，按两次计算
func ContainInt32s(array []int32, elem ...int32) bool {
	if len(elem) == 0 {
		return false
	}
	if len(elem) == 1 {
		return GetElementIndexInt32(array, elem[0]) >= 0
	}
	counter := make(map[int32]int, len(array))
	var (
		ok    bool
		count int
	)
	for _, item := range array {
		if _, ok = counter[item]; ok {
			counter[item]++
		} else {
			counter[item] = 1
		}
	}
	for _, e := range elem {
		if count, ok = counter[e]; !ok || count <= 0 {
			return false
		} else {
			counter[e]--
		}
	}
	return true
}
