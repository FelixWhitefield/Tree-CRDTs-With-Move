package main

import (
	"fmt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
)


func main() {
	l1 := clocks.NewLamport()

	l2 := clocks.NewLamport()

	fmt.Println(l1.Compare(*l2))


	v1 := clocks.NewVectorClock()
	v1.Inc()
	fmt.Println(v1)

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

	fmt.Println("Before ", v)

	v.Merge(v2.CloneTimestamp())

	fmt.Println("After ", v)

	fmt.Println(v.CompareTimestamp(v2.CloneTimestamp()))

	m := treecrdt.NewOpMove(*c, 1, 2, "hi")
	_ = m

	newf := func(a int, b int) int {
		return a + b
	}
	fmt.Println(newf(1, 2))


}
