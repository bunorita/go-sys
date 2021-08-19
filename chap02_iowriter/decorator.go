package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
)

// io.Writerのデコレータ

// 2.4.6 io.MultiWriter 複数の内部writerに渡す
func multiWriter() {
	file, err := os.Create("multiwriter.txt")
	if err != nil {
		panic(err)
	}
	w := io.MultiWriter(file, os.Stdout)
	_, err = io.WriteString(w, "io.MultiWriter example\n")
	if err != nil {
		panic(err)
	}
}

// gzip.Writer データをgzip圧縮して内部writerに渡す
func gzipWriter() {
	file, err := os.Create("test.txt.gz")
	if err != nil {
		panic(err)
	}
	w := gzip.NewWriter(file)
	w.Header.Name = "test.txt"
	_, err = io.WriteString(w, "gzip.Writer example\n")
	if err != nil {
		panic(err)
	}
	w.Close()
}

// bufio.Writer 出力結果を一時的に貯めておき、ある分量ごとに書き出す
func bufioWriter() {
	buf := bufio.NewWriter(os.Stdout)
	buf.WriteString("bufio.Writer ")
	if err := buf.Flush(); err != nil {
		panic(err)
	}
	buf.WriteString("example\n")
	if err := buf.Flush(); err != nil {
		panic(err)
	}
}
