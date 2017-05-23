package connector

import (
	"database/sql"
	"fmt"

	"file_etl_importer/load"
)

const (
	connec = "postgreSQL"
)

type Connector struct {
	Database load.Repository
	name     string
}

func (self *Connector) CreateRepository(fields []string) {

	switch self.name {
	case "postgreSQL":
		stmt, err := self.Database.BeginTransaction()

		if err != nil {
			fmt.Println("CreateRepository: ", err)
		} else {
			self.Database.CreateDatabasePostgres(stmt, fields)
		}
	case "mongoDB":
		//
	case "rabbitMQ":
		//
	}
}

func (self *Connector) SendDataToLoad(registerList []string, stmt *sql.Tx) {
	//fmt.Println("new register: ", reg)

	switch self.name {
	case "postgreSQL":
		self.Database.InsertBatch(stmt, registerList)
		//self.Database.Commit(stmt)
	case "mongoDB":
		//
	case "rabbitMQ":
		//
	}
}

func newPostgresConnector() Connector {
	repository := load.NewDatabasePostgres()

	self := Connector{
		Database: repository,
		name:     "postgreSQL",
	}

	return self
}

func NewConnector() *Connector {
	switch connec {
	case "postgreSQL":
		self := newPostgresConnector()
		return &self
	case "mongoDB":
		//
	case "rabbitMQ":
		//
	}
	return nil
}
