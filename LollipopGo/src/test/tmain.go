package main

import "fmt"

func main() {

	x := []string{"a", "b", "c"}

	for k, v := range x {
		fmt.Print(k)
		fmt.Print(v)
	}
	fmt.Println("")
	for v := range x {

		fmt.Print(v)
	}
}
