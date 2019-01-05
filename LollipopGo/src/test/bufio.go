package main

import (
	"bufio"
	"fmt"
	"os"
)

func main1() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		counts[input.Text()]++

	}

	for line, n := range counts {
		if n > 1 {
			fmt.Println(n, line)
		}
	}
}
