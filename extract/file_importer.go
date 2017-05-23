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

	var registerList []string
	batchSize := c.Processing.BatchSizeCommit

	var wg sync.WaitGroup

	for scanner.Scan() {
		if count == 0 {
			// create fields with head file
			fiel := strings.Split(scanner.Text(), ",")
			f.fieldsLenght = len(fiel)
			f.fields = fiel
			fmt.Println("file fields: ", f.fields)
			f.connec.CreateRepository(f.fields)
			count++
		} else {
			registerList = append(registerList, scanner.Text())
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

				//registerList = nil
			}
		}
	}
	//https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/
	//time.Sleep(1000000000000)

	fmt.Println("Finished. wg.Wait now...")

	wg.Wait()

	fmt.Println("Number of processed rows: ", count-1)

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
