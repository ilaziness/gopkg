package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/singleflight"
)

// singleflight 提供了重复函数只执行一次的机制
// 原理是group维护了一个key作为键的map，第一个进来的key会存到map里面，后续相同key会判断是否已存在，已存在的话会等待直到拿到第一个进来的执行结果后返回

func main() {
	g := new(singleflight.Group)
	key := "testdata"

	for n := 0; n < 10; n++ {
		// 模拟10个goroutine同时获取同样key的数据，getData只会执行一次，其他得到的都是相同的结果
		go func() {
			val, _, _ := g.Do(key, func() (interface{}, error) {
				return getData(key, n), nil
			})
			log.Println(n, "data:", val)
		}()
	}

	time.Sleep(time.Second * 2)
}

func getData(_ string, c int) string {
	time.Sleep(time.Millisecond * 50)
	log.Println(c, "get data func")
	return fmt.Sprintf("data%d", c)
}
