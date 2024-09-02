package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func (rs *sql.Rows) Close() error
// func (rs *sql.Rows) ColumnTypes() ([]*sql.ColumnType, error)
// func (rs *sql.Rows) Columns() ([]string, error)
// func (rs *sql.Rows) Err() error
// func (rs *sql.Rows) Next() bool
// func (rs *sql.Rows) NextResultSet() bool
// func (rs *sql.Rows) Scan(dest ...any) error

type Rows struct {
	CloseFake            func() error
	CountOfClose         int
	ColumnTypesFake      func() ([]*sql.ColumnType, error)
	CountOfColumnTypes   int
	ColumnsFake          func() ([]string, error)
	CountOfColumns       int
	ErrFake              func() error
	CountOfErr           int
	NextFake             func() bool
	CountOfNext          int
	NextResultSetFake    func() bool
	CountOfNextResultSet int
	ScanFake             func(dest ...any) error
	CountOfScan          int
}

func (rows *Rows) Close() error {
	rows.CountOfClose++
	return rows.CloseFake()
}

func (rows *Rows) ColumnTypes() ([]*sql.ColumnType, error) {
	rows.CountOfColumnTypes++
	return rows.ColumnTypesFake()
}

func (rows *Rows) Columns() ([]string, error) {
	rows.CountOfColumns++
	return rows.ColumnsFake()
}

func (rows *Rows) Err() error {
	rows.CountOfErr++
	return rows.ErrFake()
}

func (rows *Rows) Next() bool {
	rows.CountOfNext++
	return rows.NextFake()
}

func (rows *Rows) NextResultSet() bool {
	rows.CountOfNextResultSet++
	return rows.NextResultSetFake()
}

func (rows *Rows) Scan(dest ...any) error {
	rows.CountOfScan++
	return rows.ScanFake(dest...)
}

func (rows *Rows) VerifyCallCounts(t *testing.T, expected *Rows) {
	assert.Equal(t, expected.CountOfClose, rows.CountOfClose)
	assert.Equal(t, expected.CountOfColumnTypes, rows.CountOfColumnTypes)
	assert.Equal(t, expected.CountOfColumns, rows.CountOfColumns)
	assert.Equal(t, expected.CountOfErr, rows.CountOfErr)
	assert.Equal(t, expected.CountOfNext, rows.CountOfNext)
	assert.Equal(t, expected.CountOfNextResultSet, rows.CountOfNextResultSet)
	assert.Equal(t, expected.CountOfScan, rows.CountOfScan)
}

type Conn struct {
	PrepareFake         func(string) (*sql.Stmt, error)
	CountOfPrepare      int
	ExecContextFake     func(context.Context, string, ...any) (sql.Result, error)
	CountOfExecContext  int
	QueryContextFake    func(context.Context, string, ...any) (*sql.Rows, error)
	CountOfQueryContext int
	PingFake            func() error
	CountOfPing         int
	CloseFake           func() error
	CountOfClose        int
}

func (dbConnFake *Conn) Prepare(query string) (*sql.Stmt, error) {
	dbConnFake.CountOfPrepare++
	return dbConnFake.PrepareFake(query)
}

func (dbConnFake *Conn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	dbConnFake.CountOfExecContext++
	return dbConnFake.ExecContextFake(ctx, query, args...)
}

func (dbConnFake *Conn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	dbConnFake.CountOfQueryContext++
	return dbConnFake.QueryContextFake(ctx, query, args...)
}

func (dbConnFake *Conn) Ping() error {
	dbConnFake.CountOfPing++
	return dbConnFake.PingFake()
}

func (dbConnFake *Conn) Close() error {
	dbConnFake.CountOfClose++
	return dbConnFake.CloseFake()
}

func (dbConnFake *Conn) VerifyCallCounts(t *testing.T, expected *Conn) {
	assert.Equal(t, expected.CountOfPrepare, dbConnFake.CountOfPrepare)
	assert.Equal(t, expected.CountOfExecContext, dbConnFake.CountOfExecContext)
	assert.Equal(t, expected.CountOfQueryContext, dbConnFake.CountOfQueryContext)
	assert.Equal(t, expected.CountOfPing, dbConnFake.CountOfPing)
	assert.Equal(t, expected.CountOfClose, dbConnFake.CountOfClose)
}
