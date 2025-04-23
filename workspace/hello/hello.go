package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	// 自定义添加了reverse.Int
	fmt.Println(reverse.String("Hello"), reverse.Int(24601))
}
