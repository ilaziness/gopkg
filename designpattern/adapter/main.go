package main

import "fmt"

// 适配器模式，是把一个接口转换成另一个接口的方法

// 旧接口
type OldDbIf interface {
	SaveToDb(data string) bool
}
type OldDb struct{}

func (o *OldDb) SaveToDb(data string) bool {
	fmt.Println("old save:", data)
	return true
}

//新接口
type NewDbIf interface {
	Save(data string) (bool, error)
}

// 新旧接口适配器
type Adapter struct {
	OldDbIf
}

func (a *Adapter) Save(data string) (bool, error) {
	return a.OldDbIf.SaveToDb(data), nil
}

// 以上例子，保存数据到数据库的接口，新旧不兼容，新增一个适配层，把旧接口转成新接口
func main() {
	ad := &Adapter{&OldDb{}}
	ad.Save("test")
}
