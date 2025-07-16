package main

import "fmt"

// 模板方法模式 - 在一个方法中定义一个算法的骨架，而将一些步骤延迟到子类中

// 抽象类接口 - 数据处理模板
type DataProcessor interface {
	// 模板方法
	ProcessData()
	// 抽象方法，由具体实现类定义
	ReadData() string
	ProcessRawData(data string) string
	SaveData(data string)
	// 钩子方法，子类可以选择性重写
	ShouldCompress() bool
}

// 基础数据处理器 - 实现模板方法
type BaseDataProcessor struct {
	DataProcessor
}

func (bdp *BaseDataProcessor) ProcessData() {
	fmt.Println("=== 开始数据处理流程 ===")

	// 步骤1: 读取数据
	rawData := bdp.ReadData()
	fmt.Printf("读取到原始数据: %s\n", rawData)

	// 步骤2: 处理数据
	processedData := bdp.ProcessRawData(rawData)
	fmt.Printf("处理后数据: %s\n", processedData)

	// 步骤3: 可选的压缩步骤（钩子方法）
	finalData := processedData
	if bdp.ShouldCompress() {
		finalData = bdp.compressData(processedData)
		fmt.Printf("压缩后数据: %s\n", finalData)
	}

	// 步骤4: 保存数据
	bdp.SaveData(finalData)

	fmt.Println("=== 数据处理流程完成 ===")
}

func (bdp *BaseDataProcessor) compressData(data string) string {
	return fmt.Sprintf("compressed(%s)", data)
}

// 默认钩子方法实现
func (bdp *BaseDataProcessor) ShouldCompress() bool {
	return false
}

// 具体实现 - CSV数据处理器
type CSVDataProcessor struct {
	*BaseDataProcessor
	filename string
}

func NewCSVDataProcessor(filename string) *CSVDataProcessor {
	processor := &CSVDataProcessor{
		BaseDataProcessor: &BaseDataProcessor{},
		filename:          filename,
	}
	processor.BaseDataProcessor.DataProcessor = processor
	return processor
}

func (csv *CSVDataProcessor) ReadData() string {
	return fmt.Sprintf("CSV数据来自文件: %s", csv.filename)
}

func (csv *CSVDataProcessor) ProcessRawData(data string) string {
	return fmt.Sprintf("解析CSV格式: %s", data)
}

func (csv *CSVDataProcessor) SaveData(data string) {
	fmt.Printf("保存CSV数据到数据库: %s\n", data)
}

// 具体实现 - JSON数据处理器
type JSONDataProcessor struct {
	*BaseDataProcessor
	apiUrl string
}

func NewJSONDataProcessor(apiUrl string) *JSONDataProcessor {
	processor := &JSONDataProcessor{
		BaseDataProcessor: &BaseDataProcessor{},
		apiUrl:            apiUrl,
	}
	processor.BaseDataProcessor.DataProcessor = processor
	return processor
}

func (json *JSONDataProcessor) ReadData() string {
	return fmt.Sprintf("JSON数据来自API: %s", json.apiUrl)
}

func (json *JSONDataProcessor) ProcessRawData(data string) string {
	return fmt.Sprintf("解析JSON格式: %s", data)
}

func (json *JSONDataProcessor) SaveData(data string) {
	fmt.Printf("保存JSON数据到缓存: %s\n", data)
}

// 重写钩子方法，启用压缩
func (json *JSONDataProcessor) ShouldCompress() bool {
	return true
}

// 具体实现 - XML数据处理器
type XMLDataProcessor struct {
	*BaseDataProcessor
	source string
}

func NewXMLDataProcessor(source string) *XMLDataProcessor {
	processor := &XMLDataProcessor{
		BaseDataProcessor: &BaseDataProcessor{},
		source:            source,
	}
	processor.BaseDataProcessor.DataProcessor = processor
	return processor
}

func (xml *XMLDataProcessor) ReadData() string {
	return fmt.Sprintf("XML数据来自: %s", xml.source)
}

func (xml *XMLDataProcessor) ProcessRawData(data string) string {
	return fmt.Sprintf("解析XML格式并转换: %s", data)
}

func (xml *XMLDataProcessor) SaveData(data string) {
	fmt.Printf("保存XML数据到文件系统: %s\n", data)
}

// 另一个模板方法示例 - 饮料制作
type BeverageMaker interface {
	MakeBeverage()
	BoilWater()
	Brew()
	PourInCup()
	AddCondiments()
	CustomerWantsCondiments() bool
}

type BaseBeverageMaker struct {
	BeverageMaker
}

func (bbm *BaseBeverageMaker) MakeBeverage() {
	fmt.Println("=== 开始制作饮料 ===")
	bbm.BoilWater()
	bbm.Brew()
	bbm.PourInCup()
	if bbm.CustomerWantsCondiments() {
		bbm.AddCondiments()
	}
	fmt.Println("=== 饮料制作完成 ===")
}

func (bbm *BaseBeverageMaker) BoilWater() {
	fmt.Println("烧开水")
}

func (bbm *BaseBeverageMaker) PourInCup() {
	fmt.Println("倒入杯中")
}

// 默认钩子方法
func (bbm *BaseBeverageMaker) CustomerWantsCondiments() bool {
	return true
}

// 咖啡制作器
type CoffeeMaker struct {
	*BaseBeverageMaker
}

func NewCoffeeMaker() *CoffeeMaker {
	maker := &CoffeeMaker{
		BaseBeverageMaker: &BaseBeverageMaker{},
	}
	maker.BaseBeverageMaker.BeverageMaker = maker
	return maker
}

func (cm *CoffeeMaker) Brew() {
	fmt.Println("用沸水冲泡咖啡")
}

func (cm *CoffeeMaker) AddCondiments() {
	fmt.Println("加糖和牛奶")
}

// 茶制作器
type TeaMaker struct {
	*BaseBeverageMaker
	addLemon bool
}

func NewTeaMaker(addLemon bool) *TeaMaker {
	maker := &TeaMaker{
		BaseBeverageMaker: &BaseBeverageMaker{},
		addLemon:          addLemon,
	}
	maker.BaseBeverageMaker.BeverageMaker = maker
	return maker
}

func (tm *TeaMaker) Brew() {
	fmt.Println("用沸水浸泡茶叶")
}

func (tm *TeaMaker) AddCondiments() {
	if tm.addLemon {
		fmt.Println("加柠檬")
	} else {
		fmt.Println("加蜂蜜")
	}
}

func (tm *TeaMaker) CustomerWantsCondiments() bool {
	return tm.addLemon
}

func main() {
	fmt.Println("=== 数据处理模板方法示例 ===")

	// CSV数据处理
	csvProcessor := NewCSVDataProcessor("data.csv")
	csvProcessor.ProcessData()

	// JSON数据处理（启用压缩）
	jsonProcessor := NewJSONDataProcessor("https://api.example.com/data")
	jsonProcessor.ProcessData()

	// XML数据处理
	xmlProcessor := NewXMLDataProcessor("config.xml")
	xmlProcessor.ProcessData()

	fmt.Println("=== 饮料制作模板方法示例 ===")

	// 制作咖啡
	coffeeMaker := NewCoffeeMaker()
	coffeeMaker.MakeBeverage()

	// 制作茶（加柠檬）
	teaMakerWithLemon := NewTeaMaker(true)
	teaMakerWithLemon.MakeBeverage()

	// 制作茶（不加柠檬）
	teaMakerWithoutLemon := NewTeaMaker(false)
	teaMakerWithoutLemon.MakeBeverage()
}
