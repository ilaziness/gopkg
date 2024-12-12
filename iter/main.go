package main

import (
	"fmt"
	"iter"
)

// 1.23版新增迭代器的用法

// 作用
// 统一迭代器实现
// 对自定义类型容器统一迭代机制

func main() {
	// push 迭代器
	s := []int{1, 2, 3, 4, 5}
	for i, v := range backward(s) {
		fmt.Println(i, v)
	}

	fmt.Println("")
	for i, v := range backward2(s) {
		fmt.Println(i, v)
	}

	fmt.Println("")
	for v := range backwardOnlyValue(s) {
		fmt.Println(v)
	}

	fmt.Println("")
	// pull迭代器，需要主动调用next获取
	// Pull2 返回的第二个参数是stop，调用会会停止迭代
	next, stop := iter.Pull2(backward(s))
	defer stop()
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		fmt.Println(k, v)
	}

	fmt.Println("")
	for k, v := range f2 {
		fmt.Println(k, v)
	}
}

// 返回值就是一个迭代器
func backward[E any](s []E) func(func(int, E) bool) {
	return func(f func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !f(i, s[i]) {
				// 结束迭代
				return
			}
		}
	}
}

// backward2和backward一样
// 返回值就是一个迭代器
func backward2[E any](s []E) iter.Seq2[int, E] {
	return func(f func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !f(i, s[i]) {
				// 结束迭代
				return
			}
		}
	}
}

func backwardOnlyValue[E any](s []E) func(func(E) bool) {
	return func(f func(E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !f(s[i]) {
				return
			}
		}
	}
}

func f2(yield func(int, string) bool) {
	for i := 0; i < 10; i++ {
		if !yield(i, fmt.Sprintf("I'm %d ", i)) {
			return
		}
	}
}
