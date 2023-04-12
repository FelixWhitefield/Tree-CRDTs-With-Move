package main

import (
	"fmt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
)

func main() {
	var c *clocks.Lamport
	c = clocks.NewLamport()
	c.Inc()
	c.Inc()
	c.Inc()
	fmt.Println(c)

	v := clocks.NewVectorClock()
	v.Inc()
	v.Inc()

	v2 := clocks.NewVectorClock()
	v2.Inc()

	v.Merge(v2.CloneTimestamp())

	fmt.Println(v)

	fmt.Println(v.CompareTimestamp(v2.CloneTimestamp()))
}