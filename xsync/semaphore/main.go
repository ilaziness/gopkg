package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/sync/semaphore"
)

// 信号量实现，应用场景是限制并发数量
// 信号量是个计数器
// 如下示例：
// 总数是5，每运行一个goroutine消耗一个，消耗数量达到5之后，后续的则需要等待前面的释放计数

func main() {
	// 最大同时运行的goroutine数量
	max := int64(5)
	w := semaphore.NewWeighted(max)

	for n := 0; n < 15; n++ {
		// Acquire获取信号量，申请资源，等待队列是先进先出
		if err := w.Acquire(context.Background(), 1); err != nil {
			log.Println(err)
			break
		}
		go func() {
			// Release 释放信号量
			defer w.Release(1)
			time.Sleep(time.Second * 2)
			log.Println(2)
		}()
	}

	// 等待所有任务执行完毕，所有都指向完毕才能获取到数量5
	log.Println("end", w.Acquire(context.TODO(), max))
}
