package main

import (
	"fmt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
)

func main() {
	c := clocks.NewLamport()
	fmt.Println(c)
}