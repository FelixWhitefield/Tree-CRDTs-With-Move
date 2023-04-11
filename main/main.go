package main

import (
	"fmt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
)

func main() {
	var clock clocks.Timestamp
	clock = clocks.NewLamport(1)
	fmt.Println(clock)
}