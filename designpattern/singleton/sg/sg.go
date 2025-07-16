package sg

import "sync"

// 单例模式
// 单例结构体首字母小写，限定访问范围，再实现一个首字母大写的访问函数，相当于static方法的作用

type Message struct {
	Count int
}

// 消息池
type messagePool struct {
	pool *sync.Pool
}

// 消息池单例
// 饿汉模式，系统加载就完成了对象的初始化
var msgPool = &messagePool{
	pool: &sync.Pool{
		New: func() any {
			return &Message{}
		},
	},
}

func MsgPoolInstance() *messagePool {
	return msgPool
}

func (mp *messagePool) AddMsg(msg *Message) {
	mp.pool.Put(msg)
}

func (mp *messagePool) GetMsg() *Message {
	return mp.pool.Get().(*Message)
}

// 懒汉模式
// 需要使用时才初始化单例对象
var msgPool2 *messagePool
var once = &sync.Once{}
var mutex = &sync.RWMutex{}

// Instance
func Instance() *messagePool {
	// once 确保只执行一次初始化动作
	once.Do(func() {
		msgPool2 = &messagePool{
			pool: &sync.Pool{
				New: func() any {
					return &Message{}
				},
			},
		}
	})

	return msgPool2
}

// Instance2 双重检验锁创建对象
// 双重检验锁的方式确保只执行一次
func Instance2() *messagePool {
	//先添加读锁，防止下面读取判读和Lock()后的赋值出现冲突
	// 注释RLock()，执行go test -race 测试对比有读锁没读锁的竞态情况
	mutex.RLock()
	if msgPool2 == nil {
		mutex.RUnlock()
		mutex.Lock()
		defer mutex.Unlock()
		if msgPool2 == nil {
			msgPool2 = &messagePool{
				pool: &sync.Pool{
					New: func() any {
						return &Message{}
					},
				},
			}
		}
	} else {
		mutex.RUnlock()
	}

	return msgPool2
}
