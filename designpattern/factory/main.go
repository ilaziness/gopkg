package main

import "fmt"

type Animal interface {
	Name() string
}

type Dog struct {
	name string
}

func (d Dog) Name() string {
	return d.name
}

type Cat struct {
	name string
}

func (c Cat) Name() string {
	return c.name
}

// 工厂对象
type FactoryAnimal struct{}

func (f FactoryAnimal) Create(name string) Animal {
	switch name {
	case "Dog":
		return Dog{name: "Dog"}
	case "Cat":
		return Cat{name: "Cat"}
	}
	return nil
}

// 单独工厂
func NewDog() Dog {
	return Dog{name: "Dog"}
}
func NewCat() Cat {
	return Cat{name: "Cat"}
}

// 抽象工厂
// 抽象工厂就是把创建对象的过程用接口抽象出来,每种子类自己实现自己的工厂，没有像上面在工厂方法里面switch判断
type AnimalCreator interface {
	Create() Animal
}

func (d Dog) Create() Animal {
	return Dog{name: "Dog1"}
}

func (d Cat) Create() Animal {
	return Cat{name: "Cat1"}
}

func FactoryAnimal2(c AnimalCreator) Animal {
	return c.Create()
}

func main() {
	fmt.Println(FactoryAnimal{}.Create("Dog").Name())
	fmt.Println(FactoryAnimal{}.Create("Cat").Name())

	fmt.Println(NewDog().Name())
	fmt.Println(NewCat().Name())

	fmt.Println(FactoryAnimal2(Dog{}).Name())
	fmt.Println(FactoryAnimal2(Cat{}).Name())
}
