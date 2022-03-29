package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// stdin()
	// file()
	// conn()
	// buf()
	// limit()
	// section()
	// convertEndian()
	// inspectPNG()
	insertTextChunk()
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

// 3.5.2 binary.Read
func convertEndian() {
	data := []byte{0x0, 0x0, 0x27, 0x10} // 10000 (32bits big endian)
	var i int32

	// 実行環境のエンディアンの数値に変換する
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Println(i) // => 10000 (little endian for Intel CPU)
}

// 3.5.3
func inspectPNG() {
	file, err := os.Open("img/lenna.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	chunks := readPNGChunks(file)
	for _, chunk := range chunks {
		dumpPNGChunk(chunk) // => IHDR, sRGB, IDAT, IEND
	}
}

// A PNG file is composed of an 8-byte signature header and chunks.
// Each chunk contains 4 fields
// * (4-byte) length of data
// * (4-byte) type
// * (length bytes) data
// * (4-byte) CRC
func readPNGChunks(file *os.File) []io.Reader {
	var chunks []io.Reader

	file.Seek(8, 0) // skip signature header
	var offset int64 = 8
	for {
		var length int32 // 4-byte
		err := binary.Read(file, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		chunks = append(chunks,
			io.NewSectionReader(file, offset, int64(length)+12)) // 4+4+length+4

		// move to head of next chunk()
		// 4(type)+lenght(data)+4(CRC)
		offset, _ = file.Seek(int64(length+8), 1)
	}
	return chunks
}

func dumpPNGChunk(chunk io.Reader) {
	var length int32 // 4-byte for length
	binary.Read(chunk, binary.BigEndian, &length)
	typ := make([]byte, 4) // 4-byte for type
	chunk.Read(typ)
	fmt.Printf("chunk %q (%d bytes)\n", string(typ), length)
	if bytes.Equal(typ, []byte("tEXt")) {
		rawText := make([]byte, length)
		chunk.Read(rawText)
		fmt.Println(string(rawText))
	}
}

// 3.5.4
func insertTextChunk() {
	file, err := os.Open("img/lenna.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	newFile, err := os.Create("img/lenna2.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer newFile.Close()

	chunks := readPNGChunks(file)
	io.WriteString(newFile, "\x89PNG\r\n\x1a\n")          // write signature header
	io.Copy(newFile, chunks[0])                           // write IHDR chunk
	io.Copy(newFile, textPNGChunk("ASCII PROGRAMMING++")) // write text chunk
	for _, chunk := range chunks[1:] {                    // write rest of chunks
		io.Copy(newFile, chunk)
	}

	// inspect newf ifle
	for _, chunk := range readPNGChunks(newFile) {
		dumpPNGChunk(chunk)
	}
}

func textPNGChunk(text string) io.Reader {
	byteData := []byte(text)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int32(len(byteData))) // write length field
	buf.WriteString("tEXt")                                   // write type field
	buf.Write(byteData)                                       // write data field
	// calculate and write CRC field
	crc := crc32.NewIEEE()
	io.WriteString(crc, "tEXt")
	binary.Write(buf, binary.BigEndian, crc.Sum32())

	return buf
}
