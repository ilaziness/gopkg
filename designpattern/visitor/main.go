package main

import "fmt"

// 访问者模式 - 表示一个作用于某对象结构中的各元素的操作，它使你可以在不改变各元素的类的前提下定义作用于这些元素的新操作

// 访问者接口
type Visitor interface {
	VisitBook(book *Book)
	VisitCD(cd *CD)
	VisitDVD(dvd *DVD)
}

// 元素接口
type Element interface {
	Accept(visitor Visitor)
	GetName() string
	GetPrice() float64
}

// 具体元素 - 书籍
type Book struct {
	name   string
	price  float64
	author string
	pages  int
}

func NewBook(name string, price float64, author string, pages int) *Book {
	return &Book{
		name:   name,
		price:  price,
		author: author,
		pages:  pages,
	}
}

func (b *Book) Accept(visitor Visitor) {
	visitor.VisitBook(b)
}

func (b *Book) GetName() string {
	return b.name
}

func (b *Book) GetPrice() float64 {
	return b.price
}

func (b *Book) GetAuthor() string {
	return b.author
}

func (b *Book) GetPages() int {
	return b.pages
}

// 具体元素 - CD
type CD struct {
	name     string
	price    float64
	artist   string
	duration int // 分钟
}

func NewCD(name string, price float64, artist string, duration int) *CD {
	return &CD{
		name:     name,
		price:    price,
		artist:   artist,
		duration: duration,
	}
}

func (c *CD) Accept(visitor Visitor) {
	visitor.VisitCD(c)
}

func (c *CD) GetName() string {
	return c.name
}

func (c *CD) GetPrice() float64 {
	return c.price
}

func (c *CD) GetArtist() string {
	return c.artist
}

func (c *CD) GetDuration() int {
	return c.duration
}

// 具体元素 - DVD
type DVD struct {
	name     string
	price    float64
	director string
	runtime  int // 分钟
}

func NewDVD(name string, price float64, director string, runtime int) *DVD {
	return &DVD{
		name:     name,
		price:    price,
		director: director,
		runtime:  runtime,
	}
}

func (d *DVD) Accept(visitor Visitor) {
	visitor.VisitDVD(d)
}

func (d *DVD) GetName() string {
	return d.name
}

func (d *DVD) GetPrice() float64 {
	return d.price
}

func (d *DVD) GetDirector() string {
	return d.director
}

func (d *DVD) GetRuntime() int {
	return d.runtime
}

// 具体访问者 - 价格计算访问者
type PriceCalculatorVisitor struct {
	totalPrice float64
}

func NewPriceCalculatorVisitor() *PriceCalculatorVisitor {
	return &PriceCalculatorVisitor{totalPrice: 0}
}

func (pcv *PriceCalculatorVisitor) VisitBook(book *Book) {
	pcv.totalPrice += book.GetPrice()
	fmt.Printf("书籍 '%s' 价格: %.2f 元\n", book.GetName(), book.GetPrice())
}

func (pcv *PriceCalculatorVisitor) VisitCD(cd *CD) {
	pcv.totalPrice += cd.GetPrice()
	fmt.Printf("CD '%s' 价格: %.2f 元\n", cd.GetName(), cd.GetPrice())
}

func (pcv *PriceCalculatorVisitor) VisitDVD(dvd *DVD) {
	pcv.totalPrice += dvd.GetPrice()
	fmt.Printf("DVD '%s' 价格: %.2f 元\n", dvd.GetName(), dvd.GetPrice())
}

func (pcv *PriceCalculatorVisitor) GetTotalPrice() float64 {
	return pcv.totalPrice
}

// 具体访问者 - 详细信息显示访问者
type DetailDisplayVisitor struct{}

func NewDetailDisplayVisitor() *DetailDisplayVisitor {
	return &DetailDisplayVisitor{}
}

func (ddv *DetailDisplayVisitor) VisitBook(book *Book) {
	fmt.Printf("书籍详情:\n")
	fmt.Printf("  书名: %s\n", book.GetName())
	fmt.Printf("  作者: %s\n", book.GetAuthor())
	fmt.Printf("  页数: %d 页\n", book.GetPages())
	fmt.Printf("  价格: %.2f 元\n", book.GetPrice())
	fmt.Println()
}

func (ddv *DetailDisplayVisitor) VisitCD(cd *CD) {
	fmt.Printf("CD详情:\n")
	fmt.Printf("  专辑: %s\n", cd.GetName())
	fmt.Printf("  艺术家: %s\n", cd.GetArtist())
	fmt.Printf("  时长: %d 分钟\n", cd.GetDuration())
	fmt.Printf("  价格: %.2f 元\n", cd.GetPrice())
	fmt.Println()
}

func (ddv *DetailDisplayVisitor) VisitDVD(dvd *DVD) {
	fmt.Printf("DVD详情:\n")
	fmt.Printf("  电影: %s\n", dvd.GetName())
	fmt.Printf("  导演: %s\n", dvd.GetDirector())
	fmt.Printf("  时长: %d 分钟\n", dvd.GetRuntime())
	fmt.Printf("  价格: %.2f 元\n", dvd.GetPrice())
	fmt.Println()
}

// 具体访问者 - 折扣计算访问者
type DiscountCalculatorVisitor struct {
	totalOriginalPrice float64
	totalDiscountPrice float64
}

func NewDiscountCalculatorVisitor() *DiscountCalculatorVisitor {
	return &DiscountCalculatorVisitor{}
}

func (dcv *DiscountCalculatorVisitor) VisitBook(book *Book) {
	originalPrice := book.GetPrice()
	discountPrice := originalPrice * 0.9 // 书籍9折
	dcv.totalOriginalPrice += originalPrice
	dcv.totalDiscountPrice += discountPrice
	fmt.Printf("书籍 '%s': 原价 %.2f 元, 折后 %.2f 元 (9折)\n",
		book.GetName(), originalPrice, discountPrice)
}

func (dcv *DiscountCalculatorVisitor) VisitCD(cd *CD) {
	originalPrice := cd.GetPrice()
	discountPrice := originalPrice * 0.8 // CD 8折
	dcv.totalOriginalPrice += originalPrice
	dcv.totalDiscountPrice += discountPrice
	fmt.Printf("CD '%s': 原价 %.2f 元, 折后 %.2f 元 (8折)\n",
		cd.GetName(), originalPrice, discountPrice)
}

func (dcv *DiscountCalculatorVisitor) VisitDVD(dvd *DVD) {
	originalPrice := dvd.GetPrice()
	discountPrice := originalPrice * 0.85 // DVD 8.5折
	dcv.totalOriginalPrice += originalPrice
	dcv.totalDiscountPrice += discountPrice
	fmt.Printf("DVD '%s': 原价 %.2f 元, 折后 %.2f 元 (8.5折)\n",
		dvd.GetName(), originalPrice, discountPrice)
}

func (dcv *DiscountCalculatorVisitor) GetTotalOriginalPrice() float64 {
	return dcv.totalOriginalPrice
}

func (dcv *DiscountCalculatorVisitor) GetTotalDiscountPrice() float64 {
	return dcv.totalDiscountPrice
}

func (dcv *DiscountCalculatorVisitor) GetTotalSavings() float64 {
	return dcv.totalOriginalPrice - dcv.totalDiscountPrice
}

// 购物车 - 对象结构
type ShoppingCart struct {
	items []Element
}

func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{
		items: make([]Element, 0),
	}
}

func (sc *ShoppingCart) AddItem(item Element) {
	sc.items = append(sc.items, item)
	fmt.Printf("添加商品: %s\n", item.GetName())
}

func (sc *ShoppingCart) Accept(visitor Visitor) {
	for _, item := range sc.items {
		item.Accept(visitor)
	}
}

func (sc *ShoppingCart) GetItemCount() int {
	return len(sc.items)
}

func main() {
	fmt.Println("=== 访问者模式示例 - 购物车系统 ===")

	// 创建购物车
	cart := NewShoppingCart()

	// 添加商品
	cart.AddItem(NewBook("设计模式", 89.00, "GoF", 395))
	cart.AddItem(NewBook("重构", 79.00, "Martin Fowler", 431))
	cart.AddItem(NewCD("周杰伦精选", 45.00, "周杰伦", 65))
	cart.AddItem(NewCD("邓丽君经典", 38.00, "邓丽君", 58))
	cart.AddItem(NewDVD("肖申克的救赎", 25.00, "Frank Darabont", 142))
	cart.AddItem(NewDVD("阿甘正传", 28.00, "Robert Zemeckis", 142))

	fmt.Printf("\n购物车中共有 %d 件商品\n\n", cart.GetItemCount())

	// 使用价格计算访问者
	fmt.Println("=== 价格计算 ===")
	priceCalculator := NewPriceCalculatorVisitor()
	cart.Accept(priceCalculator)
	fmt.Printf("总价: %.2f 元\n\n", priceCalculator.GetTotalPrice())

	// 使用详细信息显示访问者
	fmt.Println("=== 商品详细信息 ===")
	detailDisplay := NewDetailDisplayVisitor()
	cart.Accept(detailDisplay)

	// 使用折扣计算访问者
	fmt.Println("=== 折扣计算 ===")
	discountCalculator := NewDiscountCalculatorVisitor()
	cart.Accept(discountCalculator)
	fmt.Printf("\n折扣汇总:\n")
	fmt.Printf("原价总计: %.2f 元\n", discountCalculator.GetTotalOriginalPrice())
	fmt.Printf("折后总计: %.2f 元\n", discountCalculator.GetTotalDiscountPrice())
	fmt.Printf("节省金额: %.2f 元\n", discountCalculator.GetTotalSavings())
}
