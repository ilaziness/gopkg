package main

import (
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

// errgroup是对sync.waitGroup的封装
// 一组任务完成，会返回第一个错误，如果有

func main() {
	g := new(errgroup.Group)
	// 设置可以同时运行的goroutine数量
	g.SetLimit(5)

	num := 10
	for i := 0; i < num; i++ {
		// 运行新的goroutine
		g.Go(func() error {
			task(i)
			return nil
		})
	}
	// g.Wait 等待所有goroutine完成
	log.Println("end", g.Wait())
}

func task(i int) {
	time.Sleep(time.Second * 5)
	log.Println(i, "run")
}
