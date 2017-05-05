package connector

import (
	"fmt"

	"file_etl_importer/load"
)

const (
	database = "postgreSQL"
)

func CreateRepository(fields []string) {

	switch database {
	case "postgreSQL":
		db := load.NewDatabasePostgres()
		load.CreateDatabasePostgres(fields, db)
	case "mongoDB":
		//
	case "rabbitMQ":
		//
	}
}

func SendDataToLoad(reg []string) {
	fmt.Println("new register: ", reg)
}
