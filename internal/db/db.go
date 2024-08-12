package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	conn *sql.DB
}

func NewDb() Db {
	return Db{}
}

func (db *Db) Connect() error {
	conn, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	db.conn = conn
	return nil
}
