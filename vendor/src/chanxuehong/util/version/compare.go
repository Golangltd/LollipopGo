package version

// Compare 比较两个 Version 的大小.
//  返回 -1 表示 a < b
//  返回 0  表示 a == b
//  返回 +1 表示 a > b
func Compare(a, b Version) int {
	switch {
	case a.Major < b.Major:
		return -1
	case a.Major > b.Major:
		return 1
	}
	switch {
	case a.Minor < b.Minor:
		return -1
	case a.Minor > b.Minor:
		return 1
	}
	switch {
	case a.Patch < b.Patch:
		return -1
	case a.Patch > b.Patch:
		return 1
	}
	return 0
}
