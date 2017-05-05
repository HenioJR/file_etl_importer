package extract

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type file struct {
	fieldsLenght int
	fields       []string
}

func ReadCsvFile() {
	f := &file{0, nil}
	read(f)
}

/*
func read() {
	filePath, _ := filepath.Abs("./test_file_10K.csv")

	log.NewLogger("file_importer").Debugf("test debug")

	data, error := ioutil.ReadFile(filePath)

	checkError(error)
	fmt.Println(string(data))
}*/

func read(f *file) {
	file, error := os.Open("/home/henio/go/src/file_etl_importer/test_file_10K.csv")
	checkError(error)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0

	for scanner.Scan() {
		if count == 0 {
			fields := scanner.Text()
			createFields(fields, f)
			count++
		} else {
			register := scanner.Text()
			createRegister(register, f)
			count++
			break
		}

	}
	fmt.Println("Number of processed rows: ", count-1)

}

func createFields(fields string, f *file) {
	fiel := strings.Split(fields, ",")

	f.fieldsLenght = len(fiel)
	f.fields = fiel

	fmt.Println("file fields length: ", f.fieldsLenght)
	fmt.Println("file fields: ", f.fields)
}

func createRegister(register string, f *file) {
	reg := strings.Split(register, ",")

	if len(reg) != f.fieldsLenght {
		//todo save the error line in a file error
		return
	} else {
		fmt.Println("new register: ", reg)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
