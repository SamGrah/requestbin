package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// working dir for all tests is module root path
	// _ "github.com/flashlabs/rootpath"

	fake "app/internal/db/test"
	"app/internal/models"
)

func testDbSetup(t *testing.T) *Db {
	db, err := NewDb("sqlite3", ":memory:")
	assert.NoError(t, err)
	err = db.Connect()
	assert.NoError(t, err)

	path := filepath.Join("./../../", "db-schema.sql")

	c, err := os.ReadFile(path)
	assert.NoError(t, err)

	sql := string(c)
	_, err = db.conn.ExecContext(context.Background(), sql)
	assert.NoError(t, err)

	return db
}

func populatedTestDbSetup(t *testing.T) *Db {
	db := testDbSetup(t)

	testDataPath := "./test/data"
	binsFile := filepath.Join(testDataPath, "bins.sql")
	requestsFile := filepath.Join(testDataPath, "requests.sql")

	for _, file := range []string{binsFile, requestsFile} {
		c, err := os.ReadFile(file)
		assert.NoError(t, err)
		sql := string(c)
		_, err = db.conn.ExecContext(context.Background(), sql)
		assert.NoError(t, err)
	}

	return db
}

func teardownTestDb(t *testing.T, db *Db) {
	err := db.conn.Close()
	assert.NoError(t, err)
	db = nil
}

func Test_Connect(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db, err := NewDb("sqlite3", ":memory:")
		assert.NoError(t, err)

		err = db.Connect()
		assert.NoError(t, err)
	})

	t.Run("error opening db", func(t *testing.T) {
		db, err := NewDb("sqlite3", ":memory:")
		assert.NoError(t, err)
		err = db.Connect()
		assert.NoError(t, err)

		db.conn.Close() // closed sqlite in memory dbs can not longer be connected to
		err = db.Connect()
		assert.Error(t, err)
	})

	t.Run("foreign keys enabled", func(t *testing.T) {
		db := testDbSetup(t)
		defer teardownTestDb(t, db)

		query := "INSERT INTO bins (bin_id, created_at, owner) VALUES ('bin-id-1', '2023-01-01 00:00:00', 'owner-1');"
		_, err := db.conn.ExecContext(context.Background(), query)
		assert.NoError(t, err)

		// attempt to violate foriegn key (bin_id)
		query = "INSERT INTO requests (timestamp, headers, body, host, method, bin) VALUES ('2023-01-01 00:00:00', 'headers', 'body', 'host', 'method', 'non-existent-bin-id');"
		_, err = db.conn.ExecContext(context.Background(), query)
		assert.Error(t, err)
	})

	t.Run("failure to enable foreign keys", func(t *testing.T) {
		db, err := NewDb("sqlite3", "file::memory:?cache=shared")
		assert.NoError(t, err)

		conn := &fake.Conn{
			ExecContextFake: func(ctx context.Context, query string, args ...any) (sql.Result, error) {
				assert.Equal(t, "PRAGMA foreign_keys=ON;", query)
				assert.Len(t, args, 0)
				return nil, errors.New("error enabling foreign keys")
			},
		}
		db.conn = conn

		err = db.Connect()
		assert.Error(t, err)

		conn.VerifyCallCounts(t, &fake.Conn{
			CountOfExecContext: 1,
		})
	})
}

func Test_InsertBin(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db := testDbSetup(t)
		defer teardownTestDb(t, db)

		err := db.InsertBin(models.Bin{
			BinId: "bin-id",
			Owner: "owner",
		})
		assert.NoError(t, err)

		query := "SELECT * FROM bins WHERE bin_id = 'bin-id';"
		rows, err := db.conn.QueryContext(context.Background(), query)
		for rows.Next() {
			var bin models.Bin
			err = rows.Scan(&bin.BinId, &bin.CreatedAt, &bin.Owner)
			assert.NoError(t, err)
			assert.Equal(t, "bin-id", bin.BinId)
			assert.NotEqual(t, "", bin.CreatedAt)
			assert.Equal(t, "owner", bin.Owner)
		}
		assert.NoError(t, err)
	})

	t.Run("error inserting bin", func(t *testing.T) {
		db := testDbSetup(t)
		defer teardownTestDb(t, db)

		err := db.conn.Close()
		assert.NoError(t, err)

		err = db.InsertBin(models.Bin{
			BinId: "bin-id",
			Owner: "owner",
		})
		assert.Error(t, err)
	})
}

func Test_InsertRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		currenTime := time.Now()
		err := db.InsertRequest(models.Request{
			RecievedAt: currenTime,
			Headers:    "headers",
			Body:       "body",
			Host:       "host",
			Method:     "method",
			Bin:        "bin-id-1",
		})
		assert.NoError(t, err)

		query := "SELECT * FROM requests WHERE bin = 'bin';"
		rows, err := db.conn.QueryContext(context.Background(), query)
		assert.NoError(t, err)
		for rows.Next() {
			var request models.Request
			err = rows.Scan(
				&request.Id,
				&request.RecievedAt,
				&request.Headers,
				&request.Body,
				&request.Host,
				&request.Method,
				&request.Bin,
			)
			assert.NoError(t, err)
			assert.Equal(t, currenTime, request.RecievedAt)
			assert.Equal(t, "headers", request.Headers)
			assert.Equal(t, "body", request.Body)
			assert.Equal(t, "host", request.Host)
			assert.Equal(t, "method", request.Method)
			assert.Equal(t, "bin", request.Bin)
		}
	})

	t.Run("error inserting request", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		currenTime := time.Now()
		err := db.InsertRequest(models.Request{
			RecievedAt: currenTime,
			Headers:    "headers",
			Body:       "body",
			Host:       "host",
			Method:     "method",
			Bin:        "non-existent-bin-id-violated-foreign-key-constraint",
		})
		assert.Error(t, err)
	})
}

func Test_GetBinContents(t *testing.T) {
	t.Run("happy path - multiple requests present in db", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		requests, err := db.GetBinContents("bin-id-1")
		assert.NoError(t, err)
		assert.Len(t, requests, 2)

		for _, request := range requests {
			assert.Equal(t, "bin-id-1", request.Bin)
		}
	})

	t.Run("happy path - no requests present in db", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		query := "INSERT INTO bins (bin_id, created_at, owner) VALUES ('new-bin', '2023-01-01 00:00:00', 'owner-1');"
		_, err := db.conn.ExecContext(context.Background(), query)
		assert.NoError(t, err)

		requests, err := db.GetBinContents("new-bin")
		assert.Nil(t, err)
		assert.Len(t, requests, 0)
	})

	t.Run("error getting bin contents", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		err := db.conn.Close()
		assert.NoError(t, err)

		requests, err := db.GetBinContents("bin-id-1")
		assert.Nil(t, requests)
		assert.Error(t, err)
	})

	t.Run("new test", func(t *testing.T) {})
}
