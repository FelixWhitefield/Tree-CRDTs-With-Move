package main

import (
	"container/list"
	"fmt"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	"github.com/google/uuid"
	"encoding/gob"
	"bytes"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	"time"
	_ "net/http/pprof"
	"net/http"
	"log"
)

// uuid.NewUUID() for version 1's
// uuid.New() for V4

type DataA struct {
	DataA string 
	Id1 int
}

type DataB struct {
	Datab string 
	Id2 int 
}

func test(i *int) {
	ni := 10
	pni := &ni
	*i = *pni
}

type Person struct {
	Name string 
	Age int32 
}

type MapKey map[int]int

func (mk MapKey) compareTo(other MapKey) bool {
	return true
}

type Rand[T comparable] struct {
	Item T
}



func main() {
	var err error

	go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

	tcpprov := connection.NewTCPProvider(2, uuid.New())
	tcpprov2 := connection.NewTCPProvider(2, uuid.New())

	go tcpprov.Listen(1111)
	go tcpprov2.Listen(1112)

	go tcpprov.Connect("localhost:1112")

	time.Sleep(1 * time.Second)

	fmt.Println("Closing for:" , tcpprov.Id)
	tcpprov.CloseAll()

	time.Sleep(3 * time.Second)

	return
	

	nmap := make(map[int]*int)
	ill := 5
	nmap[1] = &ill
	_, ok := nmap[1] 
	fmt.Println(ok)
	nmap[1] = nil 
	_, ok = nmap[1]
	fmt.Println(ok)
	delete(nmap, 1)
	_, ok = nmap[1]
	fmt.Println(ok)

	dataA := DataA{"hi", 1}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(dataA)
	if err != nil {
		panic(err)
	}


	li := list.New()
	num := 10
	ptr := &num
	li.PushBack(ptr)

	fmt.Println(*li.Front().Value.(*int))

	*ptr = 20

	fmt.Println(*li.Front().Value.(*int))

	i := 2
	test(&i)
	fmt.Println(i)

	tree := treecrdt.NewTree[string]()
	u1 := uuid.New()
	tree.Add(u1, treecrdt.NewTreeNode(treecrdt.RootUUID, "hi"))
	u2 := uuid.New()
	tree.Add(u2, treecrdt.NewTreeNode(u1, "hi2"))

	fmt.Println(tree)

	fmt.Println(tree.GetNode(u1))
	fmt.Println(tree.GetChildren(u1))

	fmt.Println(tree.IsAncestor(u2, u1))

	tree.DeleteSubTree(u1)

	fmt.Println(tree)

	l1 := clocks.NewLamport()

	l2 := clocks.NewLamport()

	fmt.Println(l1.Compare(l2))


	v1 := clocks.NewVectorClock()
	v1.Inc()
	fmt.Println(v1)

	c := clocks.NewLamport()
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

	v.Merge(v2.Timestamp())

	fmt.Println("After ", v)

	fmt.Println(v.CompareTimestamp(v2.Timestamp()))

	newf := func(a int, b int) int {
		return a + b
	}
	fmt.Println(newf(1, 2))

}
