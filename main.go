package main

import (
	"fmt"

	"file_etl_importer/extract"
)

func main() {
	fmt.Println("Starting main function...")

	extract.ReadCsvFile()
}
