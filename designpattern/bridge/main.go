package main

import "fmt"

// 桥接模式，桥接模式将继承关系转换为关联关系，从而降低了类与类之间的耦合
// 将抽象部分与它的实现部分分离，使它们都可以独立地变化，而不影响其他的模块，减少他们之间的藕合

type Color interface {
	Show()
}

type ColorBlue struct{}

func (b ColorBlue) Show() {
	fmt.Println("color blue")
}

type ColorRed struct{}

func (b ColorRed) Show() {
	fmt.Println("color red")
}

// 形状
type Shape interface {
	Display()
}

// Rectangle 长方形
type Rectangle struct {
	Color
}

func (r Rectangle) Display() {
	fmt.Println("Rectangle, color")
	r.Color.Show()
	fmt.Println()
}

// Square 正方形
type Square struct {
	Color
}

func (s *Square) Display() {
	fmt.Println("Square, color")
	s.Color.Show()
	fmt.Println()
}

// 上例，颜色和形状抽象出来，具体实现独立，互不影响
// 桥接模式本质上就是面向接口编程，可以给系统带来很好的灵活性和可扩展性。

func main() {
	r := Rectangle{ColorRed{}}
	r.Display()
	r = Rectangle{ColorBlue{}}
	r.Display()

	s := Square{ColorRed{}}
	s.Display()
	s = Square{ColorBlue{}}
	s.Display()
}
