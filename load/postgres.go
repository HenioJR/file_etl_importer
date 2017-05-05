package load

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

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
	table        string
}

func CreateDatabasePostgres(columns []string, db *Postgres) {
	columnsPrepared := strings.Join(columns, " text,")
	columnsPrepared += " text"

	stmt, err := db.BeginTransactionPostgres()
	if err != nil {
		//self.log.Warnf("[Persist] BeginTransactionPostgres", err)
		fmt.Println("[Persist] BeginTransactionPostgres", err)
	}

	stmt.Exec("DROP TABLE IF EXISTS " + db.database + "." + db.table + ";")
	_, err = stmt.Exec("CREATE TABLE " + db.database + "." + db.table + " (" + columnsPrepared + " );")

	if err != nil {
		fmt.Println("[CreateDatabasePostgres] Exec", err)
	} else {
		db.Commit(stmt)
	}

}

func NewDatabasePostgres() *Postgres {
	fmt.Println("NewDatabasePostgres")

	var user, password, port, dbname, host, db, table string
	user = "xxx"
	password = "xxx"
	port = "5433"
	dbname = "repositorio"
	host = "xxx"
	db = "testes_henio"
	table = "file_imported"

	//log := log.NewLogger("repository." + database)
	self := Postgres{
		Driver:       "postgres",
		User:         user,
		Password:     password,
		Port:         port,
		Dbname:       dbname,
		Host:         host,
		MaxOpenConns: 50,
		MaxIdleConns: 10,
		//log:          log,
		database: db,
		table:    table,
	}
	self.getDbPostgres()

	return &self
}

func (self *Postgres) getDbPostgres() {
	if self.dbPool != nil {
		return
	}

	db, err := self.getConnectionPostgres()
	if err != nil {
		//self.log.Warnf("[getDbPostgres] getConnectionPostgres", err)
		fmt.Println("[getDbPostgres] getConnectionPostgres", err)
	}
	self.dbPool = db
	self.dbErr = err
}

func (self *Postgres) getConnectionPostgres() (*sql.DB, error) {
	db, err := sql.Open(self.Driver, self.getConnectionStringPostgres())
	db.SetMaxIdleConns(self.MaxIdleConns)
	db.SetMaxOpenConns(self.MaxOpenConns)

	return db, err
}

func (self *Postgres) getConnectionStringPostgres() (connectionString string) {
	return "sslmode=disable host=" + self.Host + " database=" + self.Dbname + " port=" + self.Port + " user=" + self.User + " password=" + self.Password
}

func (self *Postgres) BeginTransactionPostgres() (*sql.Tx, error) {
	stmt, err := self.dbPool.Begin()
	if err != nil {
		//self.log.Warnf("[BeginTransactionPostgres] Begin", err)
		fmt.Println("[BeginTransactionPostgres] Begin", err)
		return nil, err
	}

	return stmt, nil
}

func (self *Postgres) InsertPostgres(stmt *sql.Tx, table string, fields []string, values []interface{}) error {

	var stringValues []string
	for i := range values {
		stringValues = append(stringValues, "$"+strconv.Itoa(i+1))
	}

	stmtPrepared, err := stmt.Prepare("INSERT INTO " + self.database + "." + table + "(" + strings.Join(fields, ",") + ") VALUES (" + strings.Join(stringValues, ",") + ");")
	if err != nil {
		//self.log.Warnf("[InsertPostgres] Prepare", err)
		fmt.Println("[InsertPostgres] Prepare", err)
		errRollback := stmt.Rollback()
		if errRollback != nil {
			//self.log.Warnf("[InsertPostgres] Rollback", errRollback)
			fmt.Println("[InsertPostgres] Rollback", errRollback)
			return errRollback
		}
		return err
	}
	_, err = stmtPrepared.Exec(values...)
	if err != nil {
		//self.log.Warnf("[InsertPostgres] Exec", err)
		fmt.Println("[InsertPostgres] Exec", err)
		errRollback := stmt.Rollback()
		if errRollback != nil {
			//self.log.Warnf("[InsertPostgres] Rollback", errRollback)
			fmt.Println("[InsertPostgres] Rollback", errRollback)
			return errRollback
		}
		return err
	}

	return nil
}

func (self *Postgres) Commit(stmt *sql.Tx) error {
	err := stmt.Commit()
	if err != nil {
		//self.log.Warnf("[Commit] Commit", err)
		fmt.Println("[Commit] Commit", err)
		errRollback := stmt.Rollback()
		if errRollback != nil {
			//self.log.Warnf("[Commit] Rollback", errRollback)
			return errRollback
		}
		return err
	}

	return nil
}
