package extract

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"file_etl_importer/config"
	"file_etl_importer/connector"
	"file_etl_importer/transform"
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
		if count == 0 {
			// create fields with head file
			line := scanner.Text()
			line = transform.RemoveQuotes(line)
			fiel := strings.Split(line, separator)
			f.fieldsLenght = len(fiel)
			f.fields = fiel
			fmt.Println("file fields: ", f.fields)
			f.connec.CreateRepository(f.fields)
			count++
		} else {
			line := scanner.Text()
			line = transform.RemoveQuotes(line)

			registerList = append(registerList, line)
			count++

			if count%batchSize == 0 {
				for maxThreads <= numberOfThreads {
					fmt.Println("sleeping...")
					time.Sleep(2 * time.Second)
				}
				wg.Add(1)
				go func() {
					fmt.Println(">>>>>>>>>>>>>>>>>>>>>> sendData <<<<<<<<<<<<<<<<<<<<<<<")
					numberOfThreads++

					start := page * batchSize
					end := start + batchSize
					page++

					fmt.Print("start: ", start)
					fmt.Println(" - end: ", end)

					regListCopy := registerList[start:end]

					stmt, _ := f.connec.Database.BeginTransaction()
					f.connec.SendDataToLoad(regListCopy, stmt)
					f.connec.Database.Commit(stmt)
					numberOfThreads--
					wg.Done()
				}()
			}
		}
	}

	wg.Wait()

	fmt.Println("Number of processed rows: ", count-1)

}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
