package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Q2.1
func formatAndWriteFile() {
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(file, "Write to file. %d %s %f\n",
		1, "is less than", 3.14)

	err = file.Close()
	if err != nil {
		panic(err)
	}
}

// Q2.2
func csvWrite() {
	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	w := csv.NewWriter(os.Stdout)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln(err)
		}
	}

	// Write any buffered data to the underlying writer(standard output).
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalln(err)
	}
}

// Q2.3
func jsonGzipResponseAndStdout() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")

		source := map[string]string{
			"Hello": "World",
		}

		// gzip.Writer => http.ResponseWriter
		zw := gzip.NewWriter(w)

		// MultiWriter => gzip.Writer, os.Stdout
		multiW := io.MultiWriter(zw, os.Stdout)

		// json.Encoder => MultiWriter
		jsonEnc := json.NewEncoder(multiW)
		jsonEnc.SetIndent("", "	")

		if err := jsonEnc.Encode(source); err != nil {
			log.Fatalln(err)
		}
		if err := zw.Flush(); err != nil {
			log.Fatalln(err)
		}
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	})

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
