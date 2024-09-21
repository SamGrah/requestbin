package db

import (
	"app/internal/models"
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DbConn interface {
	Close() error
	Prepare(query string) (*sql.Stmt, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	Ping() error
}

type Db struct {
	conn       DbConn
	driverType string
	connStr    string
}

func NewDb(driverType, connStr string) (*Db, error) {
	return &Db{
		driverType: driverType,
		connStr:    connStr,
	}, nil
}

func (db *Db) Connect() error {
	var err error
	if db.conn == nil {
		db.conn, err = sql.Open(db.driverType, db.connStr)
		if err != nil {
			return err
		}
	}

	_, err = db.conn.ExecContext(context.Background(), "PRAGMA foreign_keys=ON;")
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) InsertBin(bin models.Bin) error {
	query := "INSERT INTO bins (bin_id, created_at, owner) VALUES (?, ?, ?)"
	_, err := db.conn.ExecContext(
		context.Background(),
		query,
		bin.BinId,
		models.TimeToString(time.Now()),
		bin.Owner)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) InsertRequest(request models.Request) error {
	query := "INSERT INTO requests (timestamp, headers, body, host, method, bin) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.conn.ExecContext(
		context.Background(),
		query,
		request.RecievedAt,
		request.Headers,
		request.Body,
		request.Host,
		request.Method,
		request.Bin,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) GetBinContents(binId string) ([]models.Request, error) {
	query := "SELECT * FROM requests WHERE bin = ?"
	rows, err := db.conn.QueryContext(context.Background(), query, binId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		err := rows.Scan(
			&request.Id,
			&request.RecievedAt,
			&request.Headers,
			&request.Body,
			&request.Host,
			&request.Method,
			&request.Bin,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}
