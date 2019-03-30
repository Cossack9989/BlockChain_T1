package main

import (
	. "./core"
	"fmt"
)

func main() {
	bc := NewBlockChain()
	var s string
	for index := 0; index < 10; index++ {
		fmt.Scanln(&s)
		bc.AddBlock(s)

	}
}
