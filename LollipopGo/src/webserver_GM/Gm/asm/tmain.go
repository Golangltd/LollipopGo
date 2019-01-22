package main

import (
	"fmt"
)

func main() {

	fmt.Println("Hello, LollipopGo!")

	data := Add(10, 8)
	fmt.Println("add:", data)
}

func Add(i, j int) int {

	return i + j

}
