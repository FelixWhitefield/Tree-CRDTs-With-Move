package main

import (
	"fmt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
)

// uuid.NewUUID() for version 1's
// uuid.New() for V4

func main() {
	tree := treecrdt.NewTree[string]()
	u1 := uuid.New()
	tree.Add(u1, *treecrdt.NewTreeNode(treecrdt.RootUUID, "hi"))
	tree.Add(uuid.New(), *treecrdt.NewTreeNode(u1, "hi2"))

	fmt.Println(tree)

	fmt.Println(tree.GetNode(u1))
	fmt.Println(tree.GetChildren(u1))

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


	newf := func(a int, b int) int {
		return a + b
	}
	fmt.Println(newf(1, 2))


}
