package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	q2_3()
}

func q2_2() {
	records := [][]string{
		{"first_name", "last_name", "username"},
		{"John", "Lennon", "john"},
		{"Paul", "McCartney", "paul"},
		{"George", "Harrison", "george"},
		{"Ringo", "Starr", "ringo"},
	}
	w := csv.NewWriter(os.Stdout)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("err writing record to csv:", err)
		}
	}
	w.Flush()
}

func q2_3() {
	http.HandleFunc("/", q2_3handler)
	http.ListenAndServe(":3000", nil)
}

func q2_3handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")

	source := map[string]string{
		"Hello": "World",
	}

	zw := gzip.NewWriter(w)
	zipAndStdout := io.MultiWriter(zw, os.Stdout)
	enc := json.NewEncoder(zipAndStdout)
	enc.SetIndent("", "    ")
	enc.Encode(source)
	zw.Flush()
}
