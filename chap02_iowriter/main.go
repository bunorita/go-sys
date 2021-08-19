package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	jsonGzipResponseAndStdout()
}

// The following implement io.Writer interface

// 2.4.1 os.File ファイル出力に書く
func file() {
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}

	_, err = file.Write([]byte("os.File example\n"))
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
}

// 2.4.2 os.Stdout 標準出力に書く
func stdout() {
	_, err := os.Stdout.Write([]byte("os.Stdout example\n"))
	if err != nil {
		panic(err)
	}
}

// 2.4.3 bytes.Buffer バッファに書く
func buffer() {
	var buf bytes.Buffer
	_, err := buf.Write([]byte("bytes.Buffer example\n"))

	// WriteString()なら文字列を直接渡せる
	// _, err := buf.WriteString("bytes.Buffer example\n")
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}

// 2.4.4 strings.Builder 書き出し専用のbytes.Buffer
// 読み出しがString()のみ
func buffer2() {
	var builder strings.Builder
	_, err := builder.Write([]byte("strings.Builder example\n"))
	if err != nil {
		panic(err)
	}

	fmt.Println(builder.String())
}

// 2.4.5 net.Conn httpリクエスト送信
func netConn() {
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte("GET / HTTP/1.0\nHost: ascii.jp\n\n"))
	// _, err = io.WriteString(conn, "GET / HTTP/1.0\nHost: ascii.jp\n\n")
	if err != nil {
		panic(err)
	}

	// net.Connはio.Readerインタフェースも満たす
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		panic(err)
	}
}

// http.Request
func netConn2() {
	req, err := http.NewRequest("GET", "http://ascii.jp", nil)
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	err = req.Write(conn)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, conn)
}

// http.ResponseWriter httpレスポンス返信
func responseWriter() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "http.ResponseWriter sample\n")
		if err != nil {
			panic(err)
		}
	})
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
