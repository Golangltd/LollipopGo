package main

import (
	"fmt"
)

var (
	strcount string = "www.Golang.Ltd"
	icount   int    = 0
)

type PlayerST struct {
	UID  int
	Data map[string]*PlayerST1
}

type PlayerST1 struct {
	UID int
}

var ddd map[string]interface{}

//------------------------------------------------------------------------------

func init() {
	ddd = make(map[string]interface{})
	s := &PlayerST{
		UID: 1,
	}
	ddd["1"] = s
}

func main() {
	for k, v := range ddd {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println(v.(*PlayerST).UID)
	}
}
