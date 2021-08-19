package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// データを整形して内部writerに渡す

// 2.4.7 fmt.Fprintf
func fprintf() {
	fmt.Fprintf(os.Stdout, "Write to os.Stdout at %v", time.Now())
}

// json.Encoder
func jsonEncoder() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "	")
	err := encoder.Encode(map[string]string{
		"example": "encoding/json",
		"hello":   "world",
	})
	if err != nil {
		panic(err)
	}
}
