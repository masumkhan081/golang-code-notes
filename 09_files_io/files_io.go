// files_io.go
// Demonstrates file reading/writing, CSV, and logging in Go.
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	// Write to file
	f, err := os.Create("example.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString("Hello, Go file IO!\n")

	// Read from file
	file, err := os.Open("example.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("Read line:", scanner.Text())
	}

	// Write CSV
	csvFile, _ := os.Create("data.csv")
	writer := csv.NewWriter(csvFile)
	writer.Write([]string{"name", "age"})
	writer.Write([]string{"Alice", "30"})
	writer.Flush()
	csvFile.Close()

	// Read CSV
	csvFile, _ = os.Open("data.csv")
	reader := csv.NewReader(csvFile)
	records, _ := reader.ReadAll()
	fmt.Println("CSV records:", records)
	csvFile.Close()

	// Logging
	log.Println("This is a log message.")
}
// Documentation:
// - Use os, bufio for file IO; encoding/csv for CSV.
// - Always handle errors and close files.
// - Edge cases: file not found, permission denied, malformed CSV.
