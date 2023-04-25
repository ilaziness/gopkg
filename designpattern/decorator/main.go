package main

import "fmt"

// 装饰模式，本体对象添加新的行为功能

// Queue 队列
type Queue struct {
	Name string
}

// SendMessage 发送消息
func (q Queue) SendMessage(s string) {
	fmt.Println("send", s)
}

// Counter 计数器，发送消息的计数功能
type Counter struct {
	Count map[string]uint
}

func (c *Counter) Inc(name string) {
	if _, ok := c.Count[name]; !ok {
		c.Count[name] = 0
	}
	c.Count[name] = c.Count[name] + 1
}

func (c *Counter) ShowCount() {
	fmt.Printf("%+v\n", c.Count)
}

var CounterInstance = &Counter{Count: make(map[string]uint)}

// CounterDecorator 计数器装饰器，给Queue添加发送消息的技术功能
type CounterDecorator struct {
	queue Queue
}

func (d CounterDecorator) SendMessage(s string) {
	d.queue.SendMessage(s)
	CounterInstance.Inc(d.queue.Name)
}

// 和代理模式有点像，不通的地方是装饰器是给对象添加功能，代理模式是对对象的访问控制
func main() {
	dec := CounterDecorator{queue: Queue{Name: "queue1"}}
	dec.SendMessage("hello world")
	CounterInstance.ShowCount()

	dec.SendMessage("hello world")
	dec.SendMessage("hello world")
	CounterInstance.ShowCount()
}
