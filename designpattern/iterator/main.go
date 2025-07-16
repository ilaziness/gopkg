package main

import "fmt"

// 迭代器模式 - 提供一种方法顺序访问一个聚合对象中各个元素，而又不暴露该对象的内部表示

// 迭代器接口
type Iterator interface {
	HasNext() bool
	Next() interface{}
	Reset()
}

// 聚合接口
type Aggregate interface {
	CreateIterator() Iterator
}

// 具体聚合 - 书架
type BookShelf struct {
	books []string
}

func NewBookShelf() *BookShelf {
	return &BookShelf{books: make([]string, 0)}
}

func (bs *BookShelf) AddBook(book string) {
	bs.books = append(bs.books, book)
}

func (bs *BookShelf) GetBook(index int) string {
	if index >= 0 && index < len(bs.books) {
		return bs.books[index]
	}
	return ""
}

func (bs *BookShelf) GetLength() int {
	return len(bs.books)
}

func (bs *BookShelf) CreateIterator() Iterator {
	return NewBookShelfIterator(bs)
}

// 具体迭代器 - 书架迭代器
type BookShelfIterator struct {
	bookShelf *BookShelf
	index     int
}

func NewBookShelfIterator(bookShelf *BookShelf) *BookShelfIterator {
	return &BookShelfIterator{
		bookShelf: bookShelf,
		index:     0,
	}
}

func (bsi *BookShelfIterator) HasNext() bool {
	return bsi.index < bsi.bookShelf.GetLength()
}

func (bsi *BookShelfIterator) Next() interface{} {
	if bsi.HasNext() {
		book := bsi.bookShelf.GetBook(bsi.index)
		bsi.index++
		return book
	}
	return nil
}

func (bsi *BookShelfIterator) Reset() {
	bsi.index = 0
}

// 反向迭代器
type ReverseBookShelfIterator struct {
	bookShelf *BookShelf
	index     int
}

func NewReverseBookShelfIterator(bookShelf *BookShelf) *ReverseBookShelfIterator {
	return &ReverseBookShelfIterator{
		bookShelf: bookShelf,
		index:     bookShelf.GetLength() - 1,
	}
}

func (rbsi *ReverseBookShelfIterator) HasNext() bool {
	return rbsi.index >= 0
}

func (rbsi *ReverseBookShelfIterator) Next() interface{} {
	if rbsi.HasNext() {
		book := rbsi.bookShelf.GetBook(rbsi.index)
		rbsi.index--
		return book
	}
	return nil
}

func (rbsi *ReverseBookShelfIterator) Reset() {
	rbsi.index = rbsi.bookShelf.GetLength() - 1
}

// 扩展书架，支持多种迭代器
type ExtendedBookShelf struct {
	*BookShelf
}

func NewExtendedBookShelf() *ExtendedBookShelf {
	return &ExtendedBookShelf{
		BookShelf: NewBookShelf(),
	}
}

func (ebs *ExtendedBookShelf) CreateReverseIterator() Iterator {
	return NewReverseBookShelfIterator(ebs.BookShelf)
}

func main() {
	// 创建书架并添加书籍
	bookShelf := NewExtendedBookShelf()
	bookShelf.AddBook("设计模式")
	bookShelf.AddBook("重构")
	bookShelf.AddBook("代码整洁之道")
	bookShelf.AddBook("算法导论")

	fmt.Println("=== 正向遍历 ===")
	iterator := bookShelf.CreateIterator()
	for iterator.HasNext() {
		book := iterator.Next().(string)
		fmt.Println("书名:", book)
	}

	fmt.Println("\n=== 反向遍历 ===")
	reverseIterator := bookShelf.CreateReverseIterator()
	for reverseIterator.HasNext() {
		book := reverseIterator.Next().(string)
		fmt.Println("书名:", book)
	}

	fmt.Println("\n=== 重置迭代器后再次遍历 ===")
	iterator.Reset()
	for iterator.HasNext() {
		book := iterator.Next().(string)
		fmt.Println("书名:", book)
	}

	// 演示Go语言的range迭代器（Go 1.23+的新特性概念）
	fmt.Println("\n=== 使用函数式迭代器 ===")
	bookShelf.ForEach(func(book string) {
		fmt.Println("书名:", book)
	})
}

// 添加函数式迭代器方法
func (bs *BookShelf) ForEach(fn func(string)) {
	for _, book := range bs.books {
		fn(book)
	}
}
