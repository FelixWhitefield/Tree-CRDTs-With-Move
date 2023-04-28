package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	ti "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
)

func BenchmarkKleppmannConcurrent(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina1 := ti.NewKTree[string](c.NewTCPProvider(2, port1))
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina2 := ti.NewKTree[string](c.NewTCPProvider(2, port2))
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina3 := ti.NewKTree[string](c.NewTCPProvider(2, port3))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		time.Sleep(100 * time.Millisecond) // Time to connect

		n := 10_000
		expectedChildren := n * 3
		b.ResetTimer()

		// Insert nodes into the trees
		for i := 0; i < n; i++ {
			Lumina1.Insert(Lumina1.Root(), "Add Node")
			Lumina2.Insert(Lumina2.Root(), "Add Node")
			Lumina3.Insert(Lumina3.Root(), "Add Node")
		}

		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		startTime := time.Now()
		for time.Since(startTime) < 10*time.Second {
			children1, _ := Lumina1.GetChildren(Lumina1.Root())
			children2, _ := Lumina2.GetChildren(Lumina2.Root())
			children3, _ := Lumina3.GetChildren(Lumina3.Root())
			if len(children1) == expectedChildren && len(children2) == expectedChildren && len(children3) == expectedChildren {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		if time.Since(startTime) >= 10*time.Second {
			b.Fatal("Timeout reached")
		}
	}
}

func BenchmarkKleppmannSequential(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina1 := ti.NewKTree[string](c.NewTCPProvider(2, port1))
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina2 := ti.NewKTree[string](c.NewTCPProvider(2, port2))
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina3 := ti.NewKTree[string](c.NewTCPProvider(2, port3))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		time.Sleep(100 * time.Millisecond) // Time to connect

		n := 10_000
		expectedChildren := n * 3
		b.ResetTimer()

		// Insert nodes into the trees
		for i := 0; i < n; i++ {
			Lumina1.Insert(Lumina1.Root(), "Add Node")
		}
		for i := 0; i < n; i++ {
			Lumina2.Insert(Lumina2.Root(), "Add Node")
		}
		for i := 0; i < n; i++ {
			Lumina3.Insert(Lumina3.Root(), "Add Node")
		}
		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		startTime := time.Now()
		for time.Since(startTime) < 10*time.Second {
			children1, _ := Lumina1.GetChildren(Lumina1.Root())
			children2, _ := Lumina2.GetChildren(Lumina2.Root())
			children3, _ := Lumina3.GetChildren(Lumina3.Root())
			if len(children1) == expectedChildren && len(children2) == expectedChildren && len(children3) == expectedChildren {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		if time.Since(startTime) >= 10*time.Second {
			b.Fatal("Timeout reached")
		}
	}
}
