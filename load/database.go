package load

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"file_etl_importer/log"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Driver       string
	User         string
	Password     string
	Port         string
	Dbname       string
	Host         string
	MaxOpenConns int
	MaxIdleConns int
	log          log.Logger
	dbPool       *sql.DB
	dbErr        error
	database     string
}

func (self *Postgres) getDb() {
	if self.dbPool != nil {
		return
	}

	db, err := self.getConnection()
	if err != nil {
		self.log.Warnf("[getDb] getConnection", err)
	}
	self.dbPool = db
	self.dbErr = err
}

func (self *Postgres) getConnection() (*sql.DB, error) {
	db, err := sql.Open(self.Driver, self.getConnectionString())
	db.SetMaxIdleConns(self.MaxIdleConns)
	db.SetMaxOpenConns(self.MaxOpenConns)

	return db, err
}

func (self *Postgres) getConnectionString() (connectionString string) {
	return "sslmode=disable host=" + self.Host + " database=" + self.Dbname + " port=" + self.Port + " user=" + self.User + " password=" + self.Password
}

func (self *Postgres) CreateDatabasePostgres(stmt *sql.Tx, columns []string, tableName string) {
	columnsPrepared := strings.Join(columns, " text,")
	columnsPrepared += " text"

	stmt.Exec("DROP TABLE IF EXISTS " + self.database + "." + tableName + ";")
	_, err := stmt.Exec("CREATE TABLE " + self.database + "." + tableName + " (" + columnsPrepared + " );")

	if err != nil {
		fmt.Println("[CreateDatabasePostgres] Exec", err)
	} else {
		self.Commit(stmt)
	}

}

func (self *Postgres) BeginTransaction() (*sql.Tx, error) {
	stmt, err := self.dbPool.Begin()
	if err != nil {
		self.log.Warnf("[BeginTransaction] Begin", err)
		return nil, err
	}

	return stmt, nil
}

func (self *Postgres) Insert(stmt *sql.Tx, table string, values []string) error {

	start := time.Now()

	query := "INSERT INTO " + self.database + "." + table + " VALUES ('" + strings.Join(values, "','") + "');"

	_, err := stmt.Exec(query)

	if err != nil {
		self.log.Warnf("[Insert] ", err)
	}

	totalTime := time.Now().Sub(start)

	fmt.Println("time insert: ", totalTime)
	return nil
}

func (self *Postgres) InsertBatch(stmt *sql.Tx, table string, registerList []string) error {

	start := time.Now()

	query := "INSERT INTO " + self.database + "." + table + " VALUES "

	for i := 0; i < len(registerList); i++ {
		reg := strings.Split(registerList[i], ",")
		query += "('" + strings.Join(reg, "','") + "'),"
	}
	query = strings.TrimRight(query, ",")
	query += ";"

	_, err := stmt.Exec(query)

	if err != nil {
		self.log.Warnf("[Insert] ", err)
	}

	totalTime := time.Now().Sub(start)

	fmt.Println("time insert: ", totalTime)
	return nil
}

func (self *Postgres) Delete(stmt *sql.Tx, table string, fields []string, values []interface{}) error {
	var stringValues []string
	for i := range values {
		stringValues = append(stringValues, (fields[i] + " = $" + strconv.Itoa(i+1)))
	}

	stmtPrepared, err := stmt.Prepare("DELETE FROM " + self.database + "." + table + " WHERE " + strings.Join(stringValues, " AND ") + ";")
	if err != nil {
		self.log.Warnf("[Delete] Prepare", err)
		errRollback := stmt.Rollback()
		if errRollback != nil {
			self.log.Warnf("[Delete] Rollback", errRollback)
			return errRollback
		}
		return err
	}
	_, err = stmtPrepared.Exec(values...)
	if err != nil {
		self.log.Warnf("[Delete] Exec", err)
		errRollback := stmt.Rollback()
		if errRollback != nil {
			self.log.Warnf("[Delete] Rollback", errRollback)
			return errRollback
		}
		return err
	}

	return nil
}

func (self *Postgres) Commit(stmt *sql.Tx) error {
	err := stmt.Commit()
	if err != nil {
		self.log.Warnf("[Commit] Commit", err)
		errRollback := stmt.Rollback()
		if errRollback != nil {
			self.log.Warnf("[Commit] Rollback", errRollback)
			return errRollback
		}
		return err
	}

	return nil
}

func NewDatabasePostgres() *Postgres {
	log := log.NewLogger("postgreSQL")
	self := Postgres{
		Driver:       "postgres",
		User:         "x",
		Password:     "x",
		Port:         "x",
		Dbname:       "x",
		Host:         "x",
		MaxOpenConns: 50,
		MaxIdleConns: 10,
		database:     "testes_henio",
		log:          log,
	}
	self.getDb()

	return &self
}
