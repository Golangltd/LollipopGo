package version

import (
	"fmt"
	"strconv"
	"strings"
)

var _ fmt.Stringer = Version{}

// Version 表示一个 major.minor.patch 格式的版本号.
type Version struct {
	Major, Minor, Patch int
}

func New(major, minor, patch int) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// String 将 Version 格式化成 major.minor.patch 的格式.
func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

// Compare 比较 v 和 v2 的大小.
//  返回 -1 表示 v < v2
//  返回 0  表示 v == v2
//  返回 +1 表示 v > v2
func (v Version) Compare(v2 Version) int {
	return Compare(v, v2)
}

// Parse 解析 x, x.y, x.y.z 格式的字符串到 Version 对象, 如果成功 ok 为 true, 否则为 false.
func Parse(str string) (v Version, ok bool) {
	if str == "" {
		return
	}

	var (
		index int
		err   error
	)

	// 获取 Major
	index = strings.IndexByte(str, '.')
	switch {
	case index > 0:
		v.Major, err = strconv.Atoi(str[:index])
		if err != nil {
			return
		}
		str = str[index+1:]
		if str == "" {
			ok = true
			return
		}
	case index == 0:
		return
	case index < 0:
		v.Major, err = strconv.Atoi(str)
		if err != nil {
			return
		}
		ok = true
		return
	}

	// 获取 Minor
	index = strings.IndexByte(str, '.')
	switch {
	case index > 0:
		v.Minor, err = strconv.Atoi(str[:index])
		if err != nil {
			return
		}
		str = str[index+1:]
		if str == "" {
			ok = true
			return
		}
	case index == 0:
		return
	case index < 0:
		v.Minor, err = strconv.Atoi(str)
		if err != nil {
			return
		}
		ok = true
		return
	}

	// 获取 Patch
	v.Patch, err = strconv.Atoi(str)
	if err != nil {
		return
	}
	ok = true
	return
}
