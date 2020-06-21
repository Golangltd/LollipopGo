package tz

import (
	"time"
)

const (
	FullFormat = "2006-01-02 15:04:05" //最常用的格式
	DayFormat  = "2006-01-02"
)

func TsToDateStr(ts int64) string {
	return GetLocalStr(time.Unix(ts, 0).UTC(), "")
}
func TsToDateTimeStr(ts int64) string {
	return GetLocalStr(time.Unix(ts, 0).UTC(), FullFormat)
}

func GetTodayStr() string {
	return GetLocalStr(time.Now().UTC(), "")
}

func IndiaTimezone() *time.Location {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return loc
}

// GetLocalStr change utc time to local date str
func GetLocalStr(base time.Time, format string) string {
	if format == "" {
		format = DayFormat
	}
	return base.In(IndiaTimezone()).Format(format)
}

//如果想要将本地时间转换成UTC，直接用UTC()方法即可
//如果解析字符串，对应的是本地时间且字符串中没有时区，使用time.ParseInLocation(ChinaTimeZone())
func UTCToLocal(base time.Time) time.Time {
	return base.In(IndiaTimezone())
}

// IsSameDay check if two time is same day locally
func IsSameDay(l time.Time, r time.Time) bool {
	return GetLocalStr(l, "") == GetLocalStr(r, "")
}

func GetNowTsMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetNowTs() int64 {
	return time.Now().Unix()
}

func Schedule(what func(), delay time.Duration, stop chan bool) {
	DynamicSchedule(what, &delay, stop)
}

func LocalNow() time.Time {
	return UTCToLocal(time.Now().UTC())
}

//可以动态修改延迟时间、可关闭的定时器
func DynamicSchedule(what func(), delayAddr *time.Duration, stop chan bool) {
	go func() {
		for {
			select {
			case <-time.After(*delayAddr):
				what()
			case <-stop:
				return
			}
		}
	}()
}
