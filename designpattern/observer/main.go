package main

import "fmt"

// 观察者模式 - 定义对象间的一种一对多的依赖关系，当一个对象的状态发生改变时，所有依赖于它的对象都得到通知并被自动更新

// 观察者接口
type Observer interface {
	Update(subject Subject)
	GetID() string
}

// 主题接口
type Subject interface {
	RegisterObserver(Observer)
	RemoveObserver(Observer)
	NotifyObservers()
}

// 具体主题 - 天气数据
type WeatherData struct {
	observers   []Observer
	temperature float64
	humidity    float64
	pressure    float64
}

func NewWeatherData() *WeatherData {
	return &WeatherData{
		observers: make([]Observer, 0),
	}
}

func (wd *WeatherData) RegisterObserver(observer Observer) {
	wd.observers = append(wd.observers, observer)
	fmt.Printf("观察者 %s 已注册\n", observer.GetID())
}

func (wd *WeatherData) RemoveObserver(observer Observer) {
	for i, obs := range wd.observers {
		if obs.GetID() == observer.GetID() {
			wd.observers = append(wd.observers[:i], wd.observers[i+1:]...)
			fmt.Printf("观察者 %s 已移除\n", observer.GetID())
			break
		}
	}
}

func (wd *WeatherData) NotifyObservers() {
	for _, observer := range wd.observers {
		observer.Update(wd)
	}
}

func (wd *WeatherData) SetMeasurements(temperature, humidity, pressure float64) {
	wd.temperature = temperature
	wd.humidity = humidity
	wd.pressure = pressure
	wd.NotifyObservers()
}

func (wd *WeatherData) GetTemperature() float64 {
	return wd.temperature
}

func (wd *WeatherData) GetHumidity() float64 {
	return wd.humidity
}

func (wd *WeatherData) GetPressure() float64 {
	return wd.pressure
}

// 具体观察者 - 当前天气显示
type CurrentConditionsDisplay struct {
	id          string
	temperature float64
	humidity    float64
}

func NewCurrentConditionsDisplay(id string) *CurrentConditionsDisplay {
	return &CurrentConditionsDisplay{id: id}
}

func (ccd *CurrentConditionsDisplay) Update(subject Subject) {
	if weatherData, ok := subject.(*WeatherData); ok {
		ccd.temperature = weatherData.GetTemperature()
		ccd.humidity = weatherData.GetHumidity()
		ccd.Display()
	}
}

func (ccd *CurrentConditionsDisplay) Display() {
	fmt.Printf("[%s] 当前天气: 温度 %.1f°C, 湿度 %.1f%%\n",
		ccd.id, ccd.temperature, ccd.humidity)
}

func (ccd *CurrentConditionsDisplay) GetID() string {
	return ccd.id
}

// 具体观察者 - 统计显示
type StatisticsDisplay struct {
	id           string
	temperatures []float64
}

func NewStatisticsDisplay(id string) *StatisticsDisplay {
	return &StatisticsDisplay{
		id:           id,
		temperatures: make([]float64, 0),
	}
}

func (sd *StatisticsDisplay) Update(subject Subject) {
	if weatherData, ok := subject.(*WeatherData); ok {
		sd.temperatures = append(sd.temperatures, weatherData.GetTemperature())
		sd.Display()
	}
}

func (sd *StatisticsDisplay) Display() {
	if len(sd.temperatures) == 0 {
		return
	}

	var sum float64
	min := sd.temperatures[0]
	max := sd.temperatures[0]

	for _, temp := range sd.temperatures {
		sum += temp
		if temp < min {
			min = temp
		}
		if temp > max {
			max = temp
		}
	}

	avg := sum / float64(len(sd.temperatures))
	fmt.Printf("[%s] 温度统计: 平均 %.1f°C, 最低 %.1f°C, 最高 %.1f°C\n",
		sd.id, avg, min, max)
}

func (sd *StatisticsDisplay) GetID() string {
	return sd.id
}

// 具体观察者 - 预报显示
type ForecastDisplay struct {
	id              string
	currentPressure float64
	lastPressure    float64
}

func NewForecastDisplay(id string) *ForecastDisplay {
	return &ForecastDisplay{id: id}
}

func (fd *ForecastDisplay) Update(subject Subject) {
	if weatherData, ok := subject.(*WeatherData); ok {
		fd.lastPressure = fd.currentPressure
		fd.currentPressure = weatherData.GetPressure()
		fd.Display()
	}
}

func (fd *ForecastDisplay) Display() {
	fmt.Printf("[%s] 天气预报: ", fd.id)
	if fd.currentPressure > fd.lastPressure {
		fmt.Println("天气转好")
	} else if fd.currentPressure < fd.lastPressure {
		fmt.Println("可能下雨")
	} else {
		fmt.Println("天气稳定")
	}
}

func (fd *ForecastDisplay) GetID() string {
	return fd.id
}

func main() {
	// 创建天气数据主题
	weatherData := NewWeatherData()

	// 创建观察者
	currentDisplay := NewCurrentConditionsDisplay("当前天气显示器")
	statisticsDisplay := NewStatisticsDisplay("统计显示器")
	forecastDisplay := NewForecastDisplay("预报显示器")

	// 注册观察者
	weatherData.RegisterObserver(currentDisplay)
	weatherData.RegisterObserver(statisticsDisplay)
	weatherData.RegisterObserver(forecastDisplay)

	fmt.Println("\n=== 第一次天气更新 ===")
	weatherData.SetMeasurements(25.0, 65.0, 1013.25)

	fmt.Println("\n=== 第二次天气更新 ===")
	weatherData.SetMeasurements(27.0, 70.0, 1015.30)

	fmt.Println("\n=== 移除统计显示器 ===")
	weatherData.RemoveObserver(statisticsDisplay)

	fmt.Println("\n=== 第三次天气更新 ===")
	weatherData.SetMeasurements(23.0, 60.0, 1010.15)
}
