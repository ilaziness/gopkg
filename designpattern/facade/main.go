package main

// 外观模式，主要是为子系统提供了一个更高层次的对外统一接口。

import "fmt"

// 子系统1
type SubSystem1 struct{}

func (s *SubSystem1) Operation1() {
	fmt.Println("SubSystem1 operation")
}

// 子系统2
type SubSystem2 struct{}

func (s *SubSystem2) Operation2() {
	fmt.Println("SubSystem2 operation")
}

// 外观
type Facade struct {
	subsystem1 *SubSystem1
	subsystem2 *SubSystem2
}

func NewFacade() *Facade {
	return &Facade{
		subsystem1: &SubSystem1{},
		subsystem2: &SubSystem2{},
	}
}

func (f *Facade) Operation() {
	f.subsystem1.Operation1()
	f.subsystem2.Operation2()
}

// 定义了两个子系统 SubSystem1 和 SubSystem2，它们分别实现了自己的操作。然后我们定义了一个外观 Facade，它包含了这两个子系统的实例，并提供了一个 Operation 方法，该方法调用了子系统的操作
// Operation提供对外的统一接口，对调用者屏蔽细节
func main() {
	facade := NewFacade()
	facade.Operation()
}
