package load

import (
	"database/sql"
)

type Repository interface {
	BeginTransaction() (*sql.Tx, error)
	CreateDatabasePostgres(stmt *sql.Tx, columns []string)
	InsertBatch(stmt *sql.Tx, registerList []string) error
	Commit(stmt *sql.Tx) error
}
