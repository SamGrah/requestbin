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

		query := "INSERT INTO bins (created_at, owner) VALUES ('2023-01-01 00:00:00', 'owner-1');"
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

func Test_CreateBin(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db := testDbSetup(t)
		defer teardownTestDb(t, db)

		currentTime := time.Now()
		owner := "owner"

		id, err := db.CreateBin(models.Bin{
			CreatedAt: time.Now(),
			Owner:     "owner",
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)

		query := "SELECT * FROM bins WHERE bin_id = 1;"
		rows, err := db.conn.QueryContext(context.Background(), query)
		for rows.Next() {
			var bin models.Bin
			err = rows.Scan(&bin.BinId, &bin.CreatedAt, &bin.Owner)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), bin.BinId)
			assert.NotEqual(t, currentTime, bin.CreatedAt)
			assert.Equal(t, owner, bin.Owner)
		}
		assert.NoError(t, err)
	})

	t.Run("error inserting bin", func(t *testing.T) {
		db := testDbSetup(t)
		defer teardownTestDb(t, db)

		err := db.conn.Close()
		assert.NoError(t, err)

		_, err = db.CreateBin(models.Bin{
			Owner: "owner",
		})
		assert.Error(t, err)
	})
}

func Test_InsertRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		currentTime := time.Now()
		req := models.Request{
			RecievedAt: currentTime,
			Body:       "new-body",
			Host:       "new-host",
			RemoteAddr: "new-remoteAddr",
			RequestUri: "new-requestUri",
			Method:     "new-method",
			Bin:        1,
		}
		_ = req.SetHeaders(map[string][]string{"headers": {"header"}})
		err := db.InsertRequest(req)

		assert.NoError(t, err)

		query := "SELECT * FROM requests WHERE bin = ?"
		rows, err := db.conn.QueryContext(context.Background(), query, req.Bin)
		assert.NoError(t, err)

		var newRequest models.Request
		for rows.Next() {
			var request models.Request
			err = rows.Scan(
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
			assert.NoError(t, err)
			if request.Method == "new-method" {
				newRequest = request
			}
		}
		assert.NotNil(t, newRequest)
		assert.Equal(t, currentTime.UTC(), newRequest.RecievedAt.UTC())
		assert.Equal(t, req.Headers, newRequest.Headers)
		assert.Equal(t, req.Body, newRequest.Body)
		assert.Equal(t, req.Host, newRequest.Host)
		assert.Equal(t, req.Method, newRequest.Method)
		assert.Equal(t, req.RemoteAddr, newRequest.RemoteAddr)
		assert.Equal(t, req.RequestUri, newRequest.RequestUri)
		assert.Equal(t, req.Bin, newRequest.Bin)
	})

	t.Run("error inserting request", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		currenTime := time.Now()
		req := models.Request{
			RecievedAt: currenTime,
			Body:       "body",
			Host:       "host",
			RemoteAddr: "remoteAddr",
			RequestUri: "requestUri",
			Method:     "method",
			Bin:        9999, // bin id that does not exist
		}
		_ = req.SetHeaders(map[string][]string{"headers": {"header"}})

		err := db.InsertRequest(req)
		assert.Error(t, err)
	})
}

func Test_GetBinContents(t *testing.T) {
	t.Run("happy path - multiple requests present in db", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		binId := int64(1)
		requests, err := db.GetBinContents(binId)
		assert.NoError(t, err)
		assert.Len(t, requests, 2)

		for _, request := range requests {
			assert.Equal(t, binId, request.Bin)
		}
	})

	t.Run("happy path - new bin", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		query := "INSERT INTO bins (created_at, owner) VALUES ('2023-01-01 00:00:00', 'owner-1');"
		res, err := db.conn.ExecContext(context.Background(), query)
		assert.NoError(t, err)

		id, err := sql.Result.LastInsertId(res)
		assert.NoError(t, err)

		requests, err := db.GetBinContents(id)
		assert.Nil(t, err)
		assert.Len(t, requests, 0)
	})

	t.Run("error getting bin contents", func(t *testing.T) {
		db := populatedTestDbSetup(t)
		defer teardownTestDb(t, db)

		err := db.conn.Close()
		assert.NoError(t, err)

		requests, err := db.GetBinContents(1)
		assert.Nil(t, requests)
		assert.Error(t, err)
	})

	t.Run("new test", func(t *testing.T) {})
}
