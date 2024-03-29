package main

import "fmt"

// Fibonacci 计算斐波那契数
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	// 不能计算的数
	if n > 93 {
		return 0, fmt.Errorf("unsupported fibonacci number %d: too large", n)
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}
