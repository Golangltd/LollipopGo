package sid

import (
	"time"
)

// unix100nano returns the number of 100-nanoseconds elapsed since "1970-01-01 00:00:00 +0000 UTC".
func unix100nano(timeNow time.Time) int64 {
	return timeNow.Unix()*1e7 + int64(timeNow.Nanosecond())/100
}

// tillNext100nano spin wait till next 100-nanosecond.
func tillNext100nano(lastTimestamp int64) int64 {
	timestamp := unix100nano(time.Now())
	for timestamp <= lastTimestamp {
		timestamp = unix100nano(time.Now())
	}
	return timestamp
}
