package main

import (
	"io"
	"log"
	"strconv"
	"testing"
	"time"

	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	ti "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
	"math/rand"
)

const useOptimizations = true

// Applies a random operations to the tree
func DoRandomOperation(tree ti.Tree[string]) {
	// 60% insert, 12% remove, 28% move
	// If a remove or move errors, insert a node instead
	// This may lead to more than 60% inserts, but that's fine
	randOp := rand.Intn(100)
	if randOp < 60 {
		// Insert
		tree.Insert(tree.Root(), "Add Node")
		return
	}
	rootChildren, _ := tree.GetChildren(tree.Root())
	if randOp < 72 {
		// Remove
		if len(rootChildren) > 0 {
			err := tree.Delete(rootChildren[rand.Intn(len(rootChildren))])
			if err != nil {
				tree.Insert(tree.Root(), "Add Node")
			}
		} else {
			tree.Insert(tree.Root(), "Add Node")
		}
	} else {
		// Move
		if len(rootChildren) > 2 {
			err := tree.Move(rootChildren[rand.Intn(len(rootChildren))], rootChildren[rand.Intn(len(rootChildren))])
			if err != nil {
				tree.Insert(tree.Root(), "Add Node")
			}
		} else {
			tree.Insert(tree.Root(), "Add Node")
		}
	}
}

func BenchmarkLuminaConcurrent_10(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(10, b)
}

func BenchmarkLuminaConcurrent_100(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(100, b)
}

func BenchmarkLuminaConcurrent_1000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(1000, b)
}

func BenchmarkLuminaConcurrent_3000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(3000, b)
}

func BenchmarkLuminaConcurrent_5000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(5000, b)
}

func BenchmarkLuminaConcurrent_7500(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(7500, b)
}

func BenchmarkLuminaConcurrent_10000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(10000, b)
}

func BenchmarkLuminaConcurrent_15000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(15000, b)
}

func BenchmarkLuminaConcurrent_20000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(20000, b)
}

func BenchmarkLuminaConcurrent_25000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaConcurrent(25000, b)
}

func LuminaConcurrent(ops int, b *testing.B) {
	for iter := 0; iter < b.N; iter++ {
		b.StopTimer()
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina1 := ti.NewLTree[string](c.NewTCPProvider(2, port1), useOptimizations)
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina2 := ti.NewLTree[string](c.NewTCPProvider(2, port2), useOptimizations)
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina3 := ti.NewLTree[string](c.NewTCPProvider(2, port3), useOptimizations)
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		time.Sleep(100 * time.Millisecond) // Time to connect

		expectedApplied := uint64(ops * 3)
		b.StartTimer()

		// Insert nodes into the trees
		for i := 0; i < ops; i++ {
			DoRandomOperation(Lumina1)
			DoRandomOperation(Lumina2)
			DoRandomOperation(Lumina3)
		}

		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		startTime := time.Now()
		for time.Since(startTime) < 12*time.Second {
			opsAp1 := Lumina1.GetTotalApplied()
			opsAp2 := Lumina2.GetTotalApplied()
			opsAp3 := Lumina3.GetTotalApplied()
			if opsAp1 == expectedApplied && opsAp2 == expectedApplied && opsAp3 == expectedApplied {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		if time.Since(startTime) >= 12*time.Second {
			b.Fatalf("Timeout reached")
		}
	}
}

func BenchmarkLuminaSequential_10(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(10, b)
}

func BenchmarkLuminaSequential_100(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(100, b)
}

func BenchmarkLuminaSequential_1000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(1000, b)
}

func BenchmarkLuminaSequential_3000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(3000, b)
}

func BenchmarkLuminaSequential_5000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(5000, b)
}

func BenchmarkLuminaSequential_7500(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(7500, b)
}

func BenchmarkLuminaSequential_10000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(10000, b)
}

func BenchmarkLuminaSequential_15000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(15000, b)
}

func BenchmarkLuminaSequential_20000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(20000, b)
}

func BenchmarkLuminaSequential_25000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaSequential(25000, b)
}

func LuminaSequential(ops int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina1 := ti.NewLTree[string](c.NewTCPProvider(2, port1), useOptimizations)
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina2 := ti.NewLTree[string](c.NewTCPProvider(2, port2), useOptimizations)
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina3 := ti.NewLTree[string](c.NewTCPProvider(2, port3), useOptimizations)
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		time.Sleep(100 * time.Millisecond) // Time to connect

		expectedApplied := uint64(ops * 3)
		b.StartTimer()

		// Insert nodes into the trees
		for i := 0; i < ops; i++ {
			DoRandomOperation(Lumina1)
			//Lumina1.Insert(Lumina1.Root(), "Add Node")
		}
		for i := 0; i < ops; i++ {
			DoRandomOperation(Lumina2)
			//Lumina2.Insert(Lumina2.Root(), "Add Node")
		}
		for i := 0; i < ops; i++ {
			DoRandomOperation(Lumina3)
			//Lumina3.Insert(Lumina3.Root(), "Add Node")
		}

		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		startTime := time.Now()
		for time.Since(startTime) < 10*time.Second {
			opsAp1 := Lumina1.GetTotalApplied()
			opsAp2 := Lumina2.GetTotalApplied()
			opsAp3 := Lumina3.GetTotalApplied()
			if opsAp1 == expectedApplied && opsAp2 == expectedApplied && opsAp3 == expectedApplied {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		if time.Since(startTime) >= 10*time.Second {
			b.Fatalf("Timeout reached")
		}
	}
}

func BenchmarkLuminaAfter_10(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(10, b)
}

func BenchmarkLuminaAfter_100(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(100, b)
}

func BenchmarkLuminaAfter_500(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(500, b)
}

func BenchmarkLuminaAfter_1000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(1000, b)
}

func BenchmarkLuminaAfter_1500(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(1500, b)
}

func BenchmarkLuminaAfter_2000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(2000, b)
}

func BenchmarkLuminaAfter_3000(b *testing.B) {
	log.SetOutput(io.Discard)
	LuminaAfter(3000, b)
}

func LuminaAfter(ops int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina1 := ti.NewLTree[string](c.NewTCPProvider(2, port1), useOptimizations)
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina2 := ti.NewLTree[string](c.NewTCPProvider(2, port2), useOptimizations)
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Lumina3 := ti.NewLTree[string](c.NewTCPProvider(2, port3), useOptimizations)

		expectedApplied := uint64(ops * 3)
		b.StartTimer()

		// Insert nodes into the trees
		for i := 0; i < ops; i++ {
			DoRandomOperation(Lumina1)
			DoRandomOperation(Lumina2)
			DoRandomOperation(Lumina3)
			// Lumina1.Insert(Lumina1.Root(), "Add Node")
			// Lumina2.Insert(Lumina2.Root(), "Add Node")
			// Lumina3.Insert(Lumina3.Root(), "Add Node")
		}

		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Lumina1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		timeout := 20 * time.Second
		startTime := time.Now()
		for time.Since(startTime) < timeout {
			opsAp1 := Lumina1.GetTotalApplied()
			opsAp2 := Lumina2.GetTotalApplied()
			opsAp3 := Lumina3.GetTotalApplied()
			if opsAp1 == expectedApplied && opsAp2 == expectedApplied && opsAp3 == expectedApplied {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		children1, _ := Lumina1.GetChildren(Lumina1.Root())
		children2, _ := Lumina2.GetChildren(Lumina2.Root())
		children3, _ := Lumina3.GetChildren(Lumina3.Root())
		if time.Since(startTime) >= timeout {
			b.Fatalf("Timeout reached: %v, %v, %v", len(children1), len(children2), len(children3))
		}
	}
}
