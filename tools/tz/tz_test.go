package tz

import (
	"fmt"
	"testing"
)

func TestGetLocalStr(t *testing.T) {
	var ts = TsToDateStr(1574960492)
	fmt.Println(ts)
}

func TestTsToDateTimeStr(t *testing.T) {
	var ts = TsToDateTimeStr(1575003617)
	fmt.Println(ts)
}
