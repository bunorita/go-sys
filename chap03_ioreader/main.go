package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"net"
	"os"
)

func main() {
	// stdin()
	// file()
	// conn()
	// buf()
	// limit()
	section()
}

// 3.4.1 os.Stdin
// implements io.Reader, io.Closer
func stdin() {
	for {
		// 5bytesずつ読む
		buffer := make([]byte, 5)
		size, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		fmt.Printf("size=%d input='%s'\n", size, string(buffer))
	}
}

// 3.4.2 os.File
// implements io.Reader, io.Writer, io.Seekr, io.Closer
func file() {
	file, err := os.Open("main.go")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	io.Copy(os.Stdout, file)
}

// 3.4.3 net.Conn
// implements io.Reader, io.Writer, io.Closer
func conn() {
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write([]byte("GET / HTTP/1.0\nHost: ascii.jp\n\n"))
	if err != nil {
		log.Fatalln(err)
	}

	io.Copy(os.Stdout, conn)
}

// 3.4.4
// bytes.Buffer
func buf() {
	buf := bytes.NewBuffer([]byte("hello"))

	buf.Write([]byte(" world"))

	io.Copy(os.Stdout, buf)
}

// 3.5.1
// io.LimitReader
func limit() {
	buf := bytes.NewBuffer([]byte("hello"))

	limitR := io.LimitReader(buf, 3)

	io.Copy(os.Stdout, limitR) // => hel
}

// io.SectionReader
func section() {
	r := strings.NewReader("Example of io.SectionReader\n")

	sectionR := io.NewSectionReader(r, 14, 7)

	io.Copy(os.Stdout, sectionR) // => Section
}
