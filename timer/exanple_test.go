package LollipopGo_timer

import (
	"fmt"
)

func ExampleTimer() {
	d := NewDispatcher(10)
	// timer 1
	d.AfterFunc(1, func() {
		fmt.Println("My name is LollipopGo")
	})
	// timer 2
	t := d.AfterFunc(1, func() {
		fmt.Println("will not print")
	})
	t.Stop()
	(<-d.ChanTimer).Cb()
}