package main

import (
	"sync"
	"testing"
)

func Benchmark_main_async(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		_, _ = walk("C:\\msi", &wg, nil)
		wg.Wait()
	}
}
func Benchmark_main_sync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = walkSync("c:\\MSI", nil)
	}
}
