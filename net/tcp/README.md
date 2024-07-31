https://github.com/smallnest/1m-go-tcp-server

https://tonybai.com/2015/11/17/tcp-programming-in-golang/

https://tonybai.com/2021/07/28/classic-blocking-network-tcp-stream-protocol-parsing-practice-in-go/

https://colobu.com/2019/02/23/1m-go-tcp-connection/

协程安全：

conn一次read和write是协程安全的，业务数据包分多次write和read协程不安全。