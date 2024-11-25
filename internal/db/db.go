package db

import (
	"app/internal/models"
	"context"
	"database/sql"
	"errors"
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

func (db *Db) CreateBin(bin models.Bin) (int64, error) {
	query := "INSERT INTO bins (created_at, owner) VALUES (?, ?)"
	res, err := db.conn.ExecContext(
		context.Background(),
		query,
		models.TimeToString(time.Now()),
		bin.Owner)
	if err != nil {
		return 0, err
	}
	id, err := sql.Result.LastInsertId(res)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("no id returned")
	}

	return id, nil
}

func (db *Db) InsertRequest(request models.Request) error {
	query := "INSERT INTO requests (timestamp, headers, body, host, remoteAddr, requestUri, method, bin) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.conn.ExecContext(
		context.Background(),
		query,
		request.RecievedAt,
		request.Headers,
		request.Body,
		request.Host,
		request.RemoteAddr,
		request.RequestUri,
		request.Method,
		request.Bin,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) GetBinContents(binId int64) ([]models.Request, error) {
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
			&request.RemoteAddr,
			&request.RequestUri,
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
