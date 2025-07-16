package main

import "fmt"

// 组合模式 - 将对象组合成树形结构以表示"部分-整体"的层次结构
// 组合模式使得用户对单个对象和组合对象的使用具有一致性

// Component 组件接口
type Component interface {
	Operation() string
	Add(Component)
	Remove(Component)
	GetChild(int) Component
}

// Leaf 叶子节点
type Leaf struct {
	name string
}

func (l *Leaf) Operation() string {
	return "Leaf " + l.name
}

func (l *Leaf) Add(Component)          {}
func (l *Leaf) Remove(Component)       {}
func (l *Leaf) GetChild(int) Component { return nil }

// Composite 组合节点
type Composite struct {
	name     string
	children []Component
}

func (c *Composite) Operation() string {
	result := "Composite " + c.name + " contains: "
	for _, child := range c.children {
		result += "[" + child.Operation() + "] "
	}
	return result
}

func (c *Composite) Add(component Component) {
	c.children = append(c.children, component)
}

func (c *Composite) Remove(component Component) {
	for i, child := range c.children {
		if child == component {
			c.children = append(c.children[:i], c.children[i+1:]...)
			break
		}
	}
}

func (c *Composite) GetChild(index int) Component {
	if index >= 0 && index < len(c.children) {
		return c.children[index]
	}
	return nil
}

func main() {
	// 创建叶子节点
	leaf1 := &Leaf{name: "1"}
	leaf2 := &Leaf{name: "2"}
	leaf3 := &Leaf{name: "3"}

	// 创建组合节点
	composite1 := &Composite{name: "A"}
	composite2 := &Composite{name: "B"}

	// 构建树形结构
	composite1.Add(leaf1)
	composite1.Add(leaf2)

	composite2.Add(leaf3)
	composite2.Add(composite1)

	// 统一调用
	fmt.Println(leaf1.Operation())
	fmt.Println(composite1.Operation())
	fmt.Println(composite2.Operation())
}
