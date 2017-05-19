package load

import (
	"database/sql"
)

type Repository interface {
	BeginTransaction() (*sql.Tx, error)
	Insert(stmt *sql.Tx, table string, values []string) error
	InsertBatch(stmt *sql.Tx, table string, registerList []string) error
	Delete(stmt *sql.Tx, table string, fields []string, values []interface{}) error
	Commit(stmt *sql.Tx) error

	CreateDatabasePostgres(stmt *sql.Tx, columns []string, tableName string)
}
