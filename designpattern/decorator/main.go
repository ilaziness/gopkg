package main

import "fmt"

// 装饰模式，本体对象添加新的行为功能

// MessageSender 消息发送接口
type MessageSender interface {
	SendMessage(s string)
}

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

// CounterDecorator 计数器装饰器，给MessageSender添加发送消息的计数功能
type CounterDecorator struct {
	sender  MessageSender
	counter *Counter
}

func NewCounterDecorator(sender MessageSender) *CounterDecorator {
	return &CounterDecorator{
		sender:  sender,
		counter: &Counter{Count: make(map[string]uint)},
	}
}

func (d *CounterDecorator) SendMessage(s string) {
	d.sender.SendMessage(s)
	// 获取发送者名称用于计数
	if queue, ok := d.sender.(Queue); ok {
		d.counter.Inc(queue.Name)
	} else {
		d.counter.Inc("unknown")
	}
}

func (d *CounterDecorator) ShowCount() {
	d.counter.ShowCount()
}

// 和代理模式有点像，不同的地方是装饰器是给对象添加功能，代理模式是对对象的访问控制
func main() {
	queue := Queue{Name: "queue1"}
	dec := NewCounterDecorator(queue)

	dec.SendMessage("hello world")
	dec.ShowCount()

	dec.SendMessage("hello world")
	dec.SendMessage("hello world")
	dec.ShowCount()
}
