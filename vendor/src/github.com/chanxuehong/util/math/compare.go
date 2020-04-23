package math

func Min(x, y int64) int64 {
	if x <= y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x >= y {
		return x
	}
	return y
}
