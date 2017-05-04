package extract

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"file_etl_importer/log"
)

func ReadCsvFile() {
	fmt.Println("readCsvFile()")

	read2()
}

func read() {
	filePath, _ := filepath.Abs("./test_file_10K.csv")

	log.NewLogger("file_importer").Debugf("test debug")

	data, error := ioutil.ReadFile(filePath)

	checkError(error)
	fmt.Println(string(data))
}

func read2() {
	file, error := os.Open("/home/henio/go/src/file_etl_importer/test_file_10K.csv")
	checkError(error)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 1

	for scanner.Scan() {
		if count == 1 {
			fmt.Println(scanner.Text())
			count++
		} else {
			continue
		}

	}

}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
