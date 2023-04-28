package main

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	//"runtime"
	"time"

	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/clocks"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"

	//"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/k"
	"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt"
	//"github.com/FelixWhitefield/Tree-CRDTs-With-Move/treecrdt/maram"

	ti "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
	"github.com/google/uuid"
)

// uuid.NewUUID() for version 1's
// uuid.New() for V4

type DataA struct {
	DataA string
	Id1   int
}

type DataB struct {
	Datab string
	Id2   uint64
}

func test(i *int) {
	ni := 10
	pni := &ni
	*i = *pni
}

type Person struct {
	Name string
	Age  int32
}

type MapKey map[int]int

func (mk MapKey) compareTo(other MapKey) bool {
	return true
}

type Rand[T comparable] struct {
	Item T
}

type LargePerson struct {
	Name          string
	Age           int32
	Height        int32
	Weight        int32
	ShoeSize      int32
	NumOfChildren int32
	NumOfPets     int32
	NumOfCars     int32
	NumOfHouses   int32
	Num           int32
	Num2          string
}

func main() {
	var err error

	mtree := ti.NewLTree[string](connection.NewTCPProvider(2, 1122))
	m2tree := ti.NewLTree[string](connection.NewTCPProvider(2, 1123))
	m3tree := ti.NewLTree[string](connection.NewTCPProvider(2, 1124))

	m2tree.ConnectionProvider().Connect("localhost:1122")
	m3tree.ConnectionProvider().Connect("localhost:1122")

	
	time.Sleep(1 * time.Second)

	for i := 0; i < 10000; i++ {
		mtree.Insert(mtree.Root(), "Felix")
	}
	for i := 0; i < 10000; i++ {
		m2tree.Insert(m2tree.Root(), "Felix")
	}
	for i := 0; i < 10000; i++ {
		m3tree.Insert(m2tree.Root(), "Felix")
	}
	
	// mid, _ := mtree.Insert(mtree.Root(), "Felix")
	// mid2, _ := m2tree.Insert(m2tree.Root(), "Felix")

	// time.Sleep(1 * time.Second)

	// for i := 0; i < 10000; i++ {
	// 	m2tree.Move(mid, mid2)
	// 	mtree.Move(mid2, mid)
	// 	mtree.Move(mid, mid2)
	// }

	time.Sleep(1 * time.Second)

	nodes11, _ := mtree.GetChildren(mtree.Root())
	nodes22, _ := m2tree.GetChildren(m2tree.Root())
	nodes33, _ := m3tree.GetChildren(m3tree.Root())

	fmt.Printf("Nodes: %v and Buf: %v \n", len(nodes11), mtree.GetBufLen())
	fmt.Printf("Nodes: %v and Buf: %v \n", len(nodes22), m2tree.GetBufLen())
	fmt.Printf("Nodes: %v and Buf: %v \n", len(nodes33), m3tree.GetBufLen())
	fmt.Println(m2tree.Equals(mtree))

	// start := time.Now()
	// for i := 0; i < 10_000; i++ {
	// 	mtree.Insert(mtree.Root(), "Felix")
	// 	m2tree.Insert(m2tree.Root(), "Felix")
	// }
	// fmt.Println("Insert 1 Mil ops EACH:", time.Since(start))

	// time.Sleep(1 * time.Second)

	// nodes, _ := m2tree.GetChildren(m2tree.Root())
	// fmt.Println("Nodes under root in Tree 2:", len(nodes))
	// nodes1, _ := mtree.GetChildren(mtree.Root())
	// fmt.Println("Nodes under root in Tree 1:", len(nodes1))

	return
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var ttree ti.Tree[string]

	ktree := ti.NewKTree[string](connection.NewTCPProvider(1, 1122))
	k2tree := ti.NewKTree[string](connection.NewTCPProvider(1, 1123))

	id, _ := k2tree.Insert(k2tree.Root(), "Felix")
	k2tree.Delete(id)

	fmt.Println(k2tree.GetChildren(id))

	return

	time.Sleep(1 * time.Second)

	ttree = ktree

	// INSERT 1 million nodes
	start := time.Now()
	for i := 0; i < 10000; i++ {
		ttree.Insert(ktree.Root(), "Felix")
		k2tree.Insert(k2tree.Root(), "Felix")
	}
	// for i := 0; i < 10000; i++ {
	// 	k2tree.Insert(k2tree.Root(), "Felix")
	// }
	fmt.Println("Insert 1 Mil ops EACH:", time.Since(start))

	k2tree.ConnectionProvider().Connect("localhost:1122")

	//ktree.Insert(uuid.Nil, "Felixadadasdsasad")

	time.Sleep(20 * time.Second)

	nodes, _ := k2tree.GetChildren(k2tree.Root())
	fmt.Println("Nodes under root in Tree 2:", len(nodes))
	nodes1, _ := ktree.GetChildren(ktree.Root())
	fmt.Println("Nodes under root in Tree 1:", len(nodes1))

	return

	tcpprov := connection.NewTCPProvider(2, 1111)
	tcpprov2 := connection.NewTCPProvider(2, 1112)
	//tcpprov3 := connection.NewTCPProvider(2, uuid.New(), 1113)

	go tcpprov.Listen()
	go tcpprov2.Listen()
	//go tcpprov3.Listen()

	go tcpprov.Connect("localhost:1112")

	go tcpprov.HandleBroadcast()

	time.Sleep(1 * time.Second)

	start = time.Now()
	fmt.Println("Sending 1 Mil ops")
	for i := 0; i < 2; i++ {
		tcpprov.BroadcastChannel() <- []byte("hi")
	}
	fmt.Println("Done sending ops")
	fmt.Println("Time taken:", time.Since(start))

	time.Sleep(1 * time.Second)

	tcpprov.CloseAll()
	tcpprov2.CloseAll()

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
