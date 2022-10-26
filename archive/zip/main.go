package main

import (
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	// 一个method id就代表一种压缩算法，所以要用其他压缩算法，要定义一个对应的method id
	method_zlib uint16 = 10
	method_gzip uint16 = 11
)

func main() {
	// 生成的zip压缩文件默认是compress/flate压缩算法

	//这里注册一下用zlib,gzip压缩算法来压缩，下面才能用method id注册使用

	// zlib
	zip.RegisterCompressor(method_zlib, func(out io.Writer) (io.WriteCloser, error) {
		return zlib.NewWriterLevel(out, flate.BestCompression)
	})
	zip.RegisterDecompressor(method_zlib, func(re io.Reader) io.ReadCloser {
		read, err := zlib.NewReader(re)
		if err != nil {
			log.Fatal(err)
		}
		return read
	})

	// gzip
	zip.RegisterCompressor(method_gzip, func(out io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriterLevel(out, gzip.BestCompression)
	})
	zip.RegisterDecompressor(method_gzip, func(re io.Reader) io.ReadCloser {
		read, err := gzip.NewReader(re)
		if err != nil {
			log.Fatal(err)
		}
		return read
	})

	// 生成用zlib,gzip压缩的压缩包可以正常读取，用操作系统的资源管理器也是可以正常解压。
	write()

	read()
}

func write() {
	// Create a buffer to write our archive to.
	//buf := new(bytes.Buffer)

	// Create a new zip archive.
	//w := zip.NewWriter(buf)

	zipFile, err := os.OpenFile("text2.zip", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()
	w := zip.NewWriter(zipFile)

	//文件名最后加斜杠会创建目录
	w.Create("test/")

	//注册压缩
	w.RegisterCompressor(method_gzip, func(out io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriterLevel(out, gzip.BestCompression)
	})

	// w.RegisterCompressor(method_zlib, func(out io.Writer) (io.WriteCloser, error) {
	// 	return zlib.NewWriterLevel(out, flate.BestCompression)
	// })

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"test/readme.txt", "This archive contains some text files."},
		{"test/gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"test/todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		// 添加一个zip文件，名称要用相对路径
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := w.Create("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	fileContent, err := os.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}

	//f.Write 可以多次调用，持续往一个文件里面写数据
	_, err = f.Write(fileContent)
	if err != nil {
		log.Fatal(err)
	}
	f.Write([]byte("\nadd new line\n"))

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func read() {
	// Open a zip archive for reading.
	r, err := zip.OpenReader("text2.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// 注册解压缩算法
	r.RegisterDecompressor(method_gzip, func(re io.Reader) io.ReadCloser {
		read, err := gzip.NewReader(re)
		if err != nil {
			log.Fatal(err)
		}
		return read
	})

	// r.RegisterDecompressor(method_zlib, func(re io.Reader) io.ReadCloser {
	// 	read, err := zlib.NewReader(re)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return read
	// })

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(os.Stdout, rc)
		// CopyN如果读取的长度超过了文件总长度会返回EOF错误
		//_, err = io.CopyN(os.Stdout, rc, 68)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		rc.Close()
		fmt.Println()
	}
}
