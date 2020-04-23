package version

import (
	"fmt"
	"strconv"
	"strings"
)

var _ fmt.Stringer = Version{}

// Version 表示一个 major.minor.build.revision 格式的版本号.
type Version struct {
	Major, Minor, Build, Revision int
}

func New(major, minor, build, revision int) Version {
	return Version{
		Major:    major,
		Minor:    minor,
		Build:    build,
		Revision: revision,
	}
}

// String 将 Version 格式化成 major.minor.build.revision 的格式.
func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Build) + "." + strconv.Itoa(v.Revision)
}

// Compare 比较 v 和 v2 的大小.
//  返回 -1 表示 v < v2
//  返回 0  表示 v == v2
//  返回 +1 表示 v > v2
func (v Version) Compare(v2 Version) int {
	return Compare(v, v2)
}

// Parse 解析 x, x.y, x.y.z, x.y.z.w 格式的字符串到 Version 对象, 如果成功 ok 为 true, 否则为 false.
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

	// 获取 Build
	index = strings.IndexByte(str, '.')
	switch {
	case index > 0:
		v.Build, err = strconv.Atoi(str[:index])
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
		v.Build, err = strconv.Atoi(str)
		if err != nil {
			return
		}
		ok = true
		return
	}

	// 获取 Revision
	v.Revision, err = strconv.Atoi(str)
	if err != nil {
		return
	}
	ok = true
	return
}
