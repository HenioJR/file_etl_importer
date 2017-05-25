package main

import (
	"fmt"
	"time"

	"file_etl_importer/extract"
)

func main() {
	fmt.Println("Starting file ETL importer...")

	start := time.Now()

	extract.ReadCsvFile()

	totalTime := time.Now().Sub(start)

	fmt.Println("Execution time: ", totalTime)
}
