package main

import "testing"

// go test -benchmem -bench .

func BenchmarkString2Byte(b *testing.B) {
	s := "abc"
	for i := 0; i < b.N; i++ {
		StringToBytes(s)
	}
}

func BenchmarkString2Byte2(b *testing.B) {
	s := "abc"
	for i := 0; i < b.N; i++ {
		_ = []byte(s)
	}
}

func BenchmarkByte2String(b *testing.B) {
	s := []byte{'a', 'b', 'c'}
	for i := 0; i < b.N; i++ {
		BytesToString(s)
	}
}

func BenchmarkByte2String2(b *testing.B) {
	s := []byte{'e', 'f', 'g'}
	for i := 0; i < b.N; i++ {
		BytesToString2(s)
	}
}

func BenchmarkByte2String3(b *testing.B) {
	s := []byte{'e', 'f', 'g'}
	for i := 0; i < b.N; i++ {
		_ = string(s)
	}
}
