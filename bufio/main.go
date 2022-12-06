package main

// bufio 实现了带缓冲I/O
// 包装了io.Reader、io.Writer对象

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	str := "abcd\n123\nagcswrf\ntest hello world!你好， 世界！"
	strReader := strings.NewReader(str)

	//////////////////// Scanner
	scanner := bufio.NewScanner(strReader)
	// bufio.ScanBytes、bufio.ScanLines、bufio.ScanRunes、bufio.ScanWords这四个函数是作为Scanner.Split分割函数来使用的
	//scanner.Split(bufio.ScanBytes)
	scanner.Split(bufio.ScanLines)
	//scanner.Split(bufio.ScanRunes)
	//scanner.Split(bufio.ScanWords)
	// Scan 安装Split的分割方式，一直读取直到结尾
	for scanner.Scan() {
		fmt.Println(scanner.Bytes(), scanner.Bytes()[0] == []byte("\n")[0], scanner.Text())
	}

	///////////////////////// Reader
	// strReader.Reset("我" + str)
	// reader := bufio.NewReader(strReader)
	// r, size, _ := reader.ReadRune()
	// // size 返回缓冲区的大小
	// fmt.Println(r, size, reader.Size())
	// r, size, _ = reader.ReadRune()
	// fmt.Println(r, size)

	reader2 := bufio.NewReader(strings.NewReader("abc\ndgdg\nsdfhdsf"))
	//r, _ := reader2.ReadString([]byte("f")[0])
	r, _ := reader2.ReadString('f')
	//r, _ = reader.ReadString([]byte("\n")[0])
	fmt.Println(r)

	///////////// Writer
	f, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	wt := bufio.NewWriter(f)
	wt.WriteString("abc\ndf")
	wt.WriteString("我是谁！")
	wt.Flush()
}
