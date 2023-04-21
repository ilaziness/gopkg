package sg

import (
	"testing"
)

// go test -race -bench=. -benchmem --run=none
func BeachmarkInstanc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go Instance()
	}
}

// go test -race -bench=. -benchmem --run=none
func BeachmarkInstanc2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go Instance2()
	}
}
