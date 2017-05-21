package extract

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"file_etl_importer/connector"
)

type file struct {
	fieldsLenght int
	fields       []string
	connec       *connector.Connector
}

func ReadCsvFile() {
	connec := connector.NewConnector()
	f := &file{0, nil, connec}
	read(f)
}

func read(f *file) {
	file, error := os.Open("/home/henio/go/src/file_etl_importer/test_file_10K.csv")
	checkError(error)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0

	var registerList []string

	for scanner.Scan() {
		if count == 0 {
			fields := scanner.Text()
			createFields(fields, f)
			count++
		} else {
			registerList = append(registerList, scanner.Text())
			count++
			//break
		}
	}

	stmt, _ := f.connec.Database.BeginTransaction()
	f.connec.SendDataToLoad(registerList, stmt)
	f.connec.Database.Commit(stmt)

	fmt.Println("Number of processed rows: ", count-1)

}

func createFields(fields string, f *file) {
	fiel := strings.Split(fields, ",")

	f.fieldsLenght = len(fiel)
	f.fields = fiel

	fmt.Println("file fields: ", f.fields)

	f.connec.CreateRepository(f.fields)
}

func createRegister(register string, fieldsLenght int) []string {
	reg := strings.Split(register, ",")

	if len(reg) != fieldsLenght {
		//todo save the error line in a file error
		return nil
	} else {
		return reg
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
