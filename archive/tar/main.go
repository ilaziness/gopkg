package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// tar全称是Tape archives，是一种可以用流的方式来管理存储文件的文件格式.
// archive/tar涵盖了GNU和BSC tar工具生成的格式

func main() {
	// 写
	log.Println("write tar")

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// 写文件
	tarFile, err := os.OpenFile("test.tar", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer tarFile.Close()
	tw2 := tar.NewWriter(tarFile)

	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "his archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling license."},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		// WriteHeader 写入头信息，并且准备好接受文件内容的写入
		// Header.Size 决定了可以写入多少字节,如果文件没有完全写入将会返回一个错误
		// 调用WriteHeader会产生一个隐式的刷新（Flush）
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		// Write写入文件内容，超过Header.Size的大小，会产生一个ErrWriteTooLong的错误
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatal(err)
		}

		if err := tw2.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw2.Write([]byte(file.Body)); err != nil {
			log.Fatal(err)
		}
	}
	// Close 刷新缓存到文件并且关闭tar
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	if err := tw2.Close(); err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////////////////////////
	// 读
	log.Println("read tar")

	tr := tar.NewReader(&buf)
	for {
		// Next 读取下一个文件
		// 没有文件了会返回io.EOF错误
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}
}
