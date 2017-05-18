package extract

import (
	"bufio"
	"database/sql"
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

/*
func read() {
	filePath, _ := filepath.Abs("./test_file_10K.csv")

	log.NewLogger("file_importer").Debugf("test debug")

	data, error := ioutil.ReadFile(filePath)

	checkError(error)
	fmt.Println(string(data))
}*/

func read(f *file) {
	file, error := os.Open("/home/henio/go/src/file_etl_importer/test_file_10.csv")
	checkError(error)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0

	stmt, _ := f.connec.Database.BeginTransaction()
	for scanner.Scan() {
		if count == 0 {
			fields := scanner.Text()
			createFields(fields, f)
			count++
		} else {
			register := scanner.Text()
			createRegister(register, f, stmt)
			count++
			//break
			fmt.Println(count)
		}
	}
	f.connec.Database.Commit(stmt)

	fmt.Println("Number of processed rows: ", count-1)

}

func createFields(fields string, f *file) {
	fiel := strings.Split(fields, ",")

	f.fieldsLenght = len(fiel)
	f.fields = fiel

	fmt.Println("file fields length: ", f.fieldsLenght)
	fmt.Println("file fields: ", f.fields)

	f.connec.CreateRepository(f.fields)
}

func createRegister(register string, f *file, stmt *sql.Tx) {
	reg := strings.Split(register, ",")

	if len(reg) != f.fieldsLenght {
		//todo save the error line in a file error
	} else {
		f.connec.SendDataToLoad(reg, stmt)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
