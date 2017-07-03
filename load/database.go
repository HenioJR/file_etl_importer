package load

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"file_etl_importer/config"
	"file_etl_importer/log"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Driver        string
	User          string
	Password      string
	Port          string
	Dbname        string
	Host          string
	MaxOpenConns  int
	MaxIdleConns  int
	log           log.Logger
	dbPool        *sql.DB
	dbErr         error
	schemaOutput  string
	tableOutput   string
	fileSeparator string
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

func (self *Postgres) CreateDatabasePostgres(stmt *sql.Tx, columns []string) {
	//put quotes if have spaces in column name
	for i := 0; i < len(columns); i++ {
		a := strings.TrimSpace(columns[i])
		b := strings.Split(a, " ")
		if len(b) > 1 {
			columns[i] = "\"" + columns[i] + "\""
		}
	}

	columnsPrepared := strings.Join(columns, " text,")
	columnsPrepared += " text"
	columnsPrepared = strings.Replace(columnsPrepared, "\"\"", "\"", -1)

	stmt.Exec("DROP TABLE IF EXISTS " + self.schemaOutput + "." + self.tableOutput + ";")
	_, err := stmt.Exec("CREATE TABLE " + self.schemaOutput + "." + self.tableOutput + " (" + columnsPrepared + " );")

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

func (self *Postgres) InsertBatch(stmt *sql.Tx, registerList []string) error {

	start := time.Now()

	query := "INSERT INTO " + self.schemaOutput + "." + self.tableOutput + " VALUES "

	for i := 0; i < len(registerList); i++ {
		if registerList[i] != "" {
			reg := strings.Split(registerList[i], self.fileSeparator)
			query += "('" + strings.Join(reg, "','") + "'),"
		}

	}
	query = strings.TrimRight(query, ",")
	query += ";"

	_, err := stmt.Exec(query)

	if err != nil {
		self.log.Warnf("[Insert] ", err)
		//self.log.Warnf("Query: ", query)
	}

	totalTime := time.Now().Sub(start)

	fmt.Println("time batch insert: ", totalTime)
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
	c := config.GetConfig()

	log := log.NewLogger("postgreSQL")
	self := Postgres{
		Driver:        c.Database.Postgres.Driver,
		User:          c.Database.Postgres.User,
		Password:      c.Database.Postgres.Password,
		Port:          c.Database.Postgres.Port,
		Dbname:        c.Database.Postgres.Dbname,
		Host:          c.Database.Postgres.Host,
		MaxOpenConns:  c.Database.Postgres.MaxOpenConns,
		MaxIdleConns:  c.Database.Postgres.MaxIdleConns,
		schemaOutput:  c.Database.Postgres.SchemaOutput,
		tableOutput:   c.Database.Postgres.TableOutput,
		log:           log,
		fileSeparator: c.File.Separator,
	}

	self.getDb()

	return &self
}
