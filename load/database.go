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
	schemaOutput string
	tableOutput  string
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
	columnsPrepared := strings.Join(columns, " text,")
	columnsPrepared += " text"

	stmt.Exec("DROP TABLE IF EXISTS " + self.schemaOutput + "." + self.tableOutput + ";")
	_, err := stmt.Exec("CREATE TABLE " + self.schemaOutput + "." + self.tableOutput + " (" + columnsPrepared + " );")

	if err != nil {
		fmt.Println("[CreateDatabasePostgres] Exec", err)
	} else {
		self.Commit(stmt)
	}

}

func (self *Postgres) BeginTransaction() (*sql.Tx, error) {

	fmt.Println("BeginTransaction")

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
			reg := strings.Split(registerList[i], ",")
			query += "('" + strings.Join(reg, "','") + "'),"
		}

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
		Driver:       c.Database.Driver,
		User:         c.Database.User,
		Password:     c.Database.Password,
		Port:         c.Database.Port,
		Dbname:       c.Database.Dbname,
		Host:         c.Database.Host,
		MaxOpenConns: c.Database.MaxOpenConns,
		MaxIdleConns: c.Database.MaxIdleConns,
		schemaOutput: c.Database.SchemaOutput,
		tableOutput:  c.Database.TableOutput,
		log:          log,
	}

	self.getDb()

	return &self
}
