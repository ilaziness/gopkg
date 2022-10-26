package main

import (
	"archive/zip"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
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

	//压缩
	/*
		w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriterLevel(out, gzip.BestCompression)
		})
	*/
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return zlib.NewWriterLevel(out, flate.BestCompression)
	})

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
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
	/*
		r.RegisterDecompressor(zip.Deflate, func(re io.Reader) io.ReadCloser {
			read, err := gzip.NewReader(re)
			if err != nil {
				log.Fatal(err)
			}
			return read
		})
	*/
	r.RegisterDecompressor(zip.Deflate, func(re io.Reader) io.ReadCloser {
		read, err := zlib.NewReader(re)
		if err != nil {
			log.Fatal(err)
		}
		return read
	})

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(os.Stdout, rc)
		//_, err = io.CopyN(os.Stdout, rc, 68)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		rc.Close()
		fmt.Println()
	}
}
