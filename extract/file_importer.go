package extract

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/file_etl_importer/config"
	"github.com/file_etl_importer/connector"
	"github.com/file_etl_importer/transform"
)

type file struct {
	fieldsLenght int
	fields       []string
	connec       *connector.Connector
}

func ReadCsvFile() {
	connec := connector.NewConnector()
	f := file{0, nil, connec}
	read(f)
}

func read(f file) {
	c := config.GetConfig()

	file, error := os.Open(c.File.PathName)
	checkError(error)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0
	page := 0
	numberOfThreads := 0
	maxThreads := c.Processing.NumberOfThreads
	batchSize := c.Processing.BatchSizeCommit
	separator := c.File.Separator

	var registerList []string
	var wg sync.WaitGroup

	for scanner.Scan() {
		line := scanner.Text()

		line = transform.RemoveQuotesInsideString(line)
		line = transform.RemoveParenthesis(line)
		//line = transform.RemoveSeparatorInsideString(line, separator)

		if c.Processing.RemoveDoubleQuote {
			line = transform.RemoveDoubleQuote(line)
		}

		if count == 0 {
			// create fields with head file
			line = transform.RemoveSpecialCharactersFromHeader(line)
			fiel := strings.Split(line, separator)
			f.fieldsLenght = len(fiel)
			f.fields = fiel
			fmt.Println("file fields: ", f.fields)
			f.connec.CreateRepository(f.fields)
			count++
		} else {
			// process data lines
			if c.Processing.ValidateLineByLine && len(strings.Split(line, separator)) != f.fieldsLenght {
				fmt.Println("Line data size error: ", line)
			} else {
				registerList = append(registerList, line)
				count++

				if count%batchSize == 0 {
					for maxThreads <= numberOfThreads {
						fmt.Println("sleeping...")
						time.Sleep(2 * time.Second)
					}
					wg.Add(1)
					go func() {
						numberOfThreads++

						start := page * batchSize
						end := start + batchSize
						page++

						saveBatch(registerList, f, start, end)

						numberOfThreads--
						wg.Done()
					}()
				}
			}
		}
	}
	wg.Wait()

	// save last incomplete batch, if necessary
	start := page * batchSize
	end := len(registerList)
	if start < len(registerList) {
		saveBatch(registerList, f, start, end)
	}

	fmt.Println("Number of processed rows: ", count-1)

}

func saveBatch(registerList []string, f file, start int, end int) {

	fmt.Print("sendData --- ")
	fmt.Print("start: ", start)
	fmt.Println(" - end: ", end)

	regListCopy := registerList[start:end]

	stmt, _ := f.connec.Database.BeginTransaction()
	f.connec.SendDataToLoad(regListCopy, stmt)
	f.connec.Database.Commit(stmt)

}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
