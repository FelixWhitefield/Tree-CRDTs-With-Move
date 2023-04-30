package main

import (
	"io"
	"log"
	"strconv"
	"testing"
	"time"

	c "github.com/FelixWhitefield/Tree-CRDTs-With-Move/connection"
	ti "github.com/FelixWhitefield/Tree-CRDTs-With-Move/treeinterface"
)

func BenchmarkKleppmannConcurrent_10(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(10, b)
}

func BenchmarkKleppmannConcurrent_100(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(100, b)
}

func BenchmarkKleppmannConcurrent_1000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(1000, b)
}

func BenchmarkKleppmannConcurrent_3000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(3000, b)
}

func BenchmarkKleppmannConcurrent_5000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(5000, b)
}

func BenchmarkKleppmannConcurrent_7500(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(7500, b)
}

func BenchmarkKleppmannConcurrent_10000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(10000, b)
}

func BenchmarkKleppmannConcurrent_15000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(15000, b)
}

func BenchmarkKleppmannConcurrent_20000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(20000, b)
}

func BenchmarkKleppmannConcurrent_25000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannConcurrent(25000, b)
}

func KleppmannConcurrent(ops int, b *testing.B) {
	for iter := 0; iter < b.N; iter++ {
		b.StopTimer()
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp1 := ti.NewKTree[string](c.NewTCPProvider(2, port1))
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp2 := ti.NewKTree[string](c.NewTCPProvider(2, port2))
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp3 := ti.NewKTree[string](c.NewTCPProvider(2, port3))
		Klepp1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Klepp1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		time.Sleep(100 * time.Millisecond) // Time to connect

		expectedApplied := uint64(ops * 3)
		b.StartTimer()

		// Insert nodes into the trees
		for i := 0; i < ops; i++ {
			DoRandomOperation(Klepp1)
			DoRandomOperation(Klepp2)
			DoRandomOperation(Klepp3)
		}

		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		startTime := time.Now()
		for time.Since(startTime) < 12*time.Second {
			opsAp1 := Klepp1.GetTotalApplied()
			opsAp2 := Klepp2.GetTotalApplied()
			opsAp3 := Klepp3.GetTotalApplied()
			if opsAp1 == expectedApplied && opsAp2 == expectedApplied && opsAp3 == expectedApplied {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		if time.Since(startTime) >= 12*time.Second {
			b.Fatalf("Timeout reached, %v", Klepp1.GetTotalApplied())
		}
	}
}
func BenchmarkKleppmannSequential_10(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(10, b)
}

func BenchmarkKleppmannSequential_100(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(100, b)
}

func BenchmarkKleppmannSequential_1000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(1000, b)
}

func BenchmarkKleppmannSequential_3000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(3000, b)
}

func BenchmarkKleppmannSequential_5000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(5000, b)
}

func BenchmarkKleppmannSequential_7500(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(7500, b)
}

func BenchmarkKleppmannSequential_10000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(10000, b)
}

func BenchmarkKleppmannSequential_15000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(15000, b)
}

func BenchmarkKleppmannSequential_20000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(20000, b)
}

func BenchmarkKleppmannSequential_25000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannSequential(25000, b)
}

func KleppmannSequential(ops int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp1 := ti.NewKTree[string](c.NewTCPProvider(2, port1))
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp2 := ti.NewKTree[string](c.NewTCPProvider(2, port2))
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp := ti.NewKTree[string](c.NewTCPProvider(2, port3))
		Klepp1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Klepp1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		time.Sleep(100 * time.Millisecond) // Time to connect

		expectedApplied := uint64(ops * 3)
		b.StartTimer()

		// Insert nodes into the trees
		for i := 0; i < ops; i++ {
			DoRandomOperation(Klepp1)
			//Klepp1.Insert(Klepp1.Root(), "Add Node")
		}
		for i := 0; i < ops; i++ {
			DoRandomOperation(Klepp2)
			//Klepp2.Insert(Klepp2.Root(), "Add Node")
		}
		for i := 0; i < ops; i++ {
			DoRandomOperation(Klepp)
			//Klepp.Insert(Klepp.Root(), "Add Node")
		}
		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		startTime := time.Now()
		for time.Since(startTime) < 10*time.Second {
			opsAp1 := Klepp1.GetTotalApplied()
			opsAp2 := Klepp2.GetTotalApplied()
			opsAp3 := Klepp.GetTotalApplied()
			if opsAp1 == expectedApplied && opsAp2 == expectedApplied && opsAp3 == expectedApplied {
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

func BenchmarkKleppmannAfter_10(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(10, b)
}

func BenchmarkKleppmannAfter_100(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(100, b)
}

func BenchmarkKleppmannAfter_500(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(500, b)
}

func BenchmarkKleppmannAfter_1000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(1000, b)
}

func BenchmarkKleppmannAfter_1500(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(1500, b)
}

func BenchmarkKleppmannAfter_2000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(2000, b)
}

func BenchmarkKleppmannAfter_3000(b *testing.B) {
	log.SetOutput(io.Discard)
	KleppmannAfter(3000, b)
}


func KleppmannAfter(ops int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		port1, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp1 := ti.NewKTree[string](c.NewTCPProvider(2, port1))
		port2, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp2 := ti.NewKTree[string](c.NewTCPProvider(2, port2))
		port3, err := GetFreePort()
		if err != nil {
			b.Fatal(err)
		}
		Klepp3 := ti.NewKTree[string](c.NewTCPProvider(2, port3))

		expectedApplied := uint64(ops * 3)
		b.StartTimer()

		// Insert nodes into the trees
		for i := 0; i < ops; i++ {
			DoRandomOperation(Klepp1)
			DoRandomOperation(Klepp2)
			DoRandomOperation(Klepp3)
			// Klepp1.Insert(Klepp1.Root(), "Add Node")
			// Klepp2.Insert(Klepp2.Root(), "Add Node")
			// Klepp3.Insert(Klepp3.Root(), "Add Node")
		}

		Klepp1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port2))
		Klepp1.ConnectionProvider().Connect("localhost:" + strconv.Itoa(port3))

		// Keep calling GetChildren() until the expected result is returned or the timeout is reached
		timeout := 20 * time.Second
		startTime := time.Now()
		for time.Since(startTime) < timeout {
			opsAp1 := Klepp1.GetTotalApplied()
			opsAp2 := Klepp2.GetTotalApplied()
			opsAp3 := Klepp3.GetTotalApplied()
			if opsAp1 == expectedApplied && opsAp2 == expectedApplied && opsAp3 == expectedApplied {
				// Expected result found, exit the loop
				break
			}
			time.Sleep(10 * time.Millisecond) // wait before calling GetChildren() again
		}
		if time.Since(startTime) >= timeout {
			b.Fatal("Timeout reached")
		}
	}
}
