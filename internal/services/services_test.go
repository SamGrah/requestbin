package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	fake "app/internal/db/test"
	"app/internal/models"
)

func Test_CreateNewBin(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var generatedBinId string
		db := fake.Db{
			InsertBinFake: func(bin models.Bin) error {
				generatedBinId = bin.BinId
				_, err := uuid.Parse(generatedBinId)
				assert.NoError(t, err)
				return nil
			},
		}

		services := New(&Deps{
			Db: &db,
		})

		binId, err := services.CreateNewBin()
		assert.NoError(t, err)
		assert.Equal(t, generatedBinId, binId)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfInsertBin: 1,
		})
	})
	t.Run("error inserting bin", func(t *testing.T) {
		db := fake.Db{
			InsertBinFake: func(bin models.Bin) error {
				return assert.AnError
			},
		}
		services := New(&Deps{
			Db: &db,
		})

		_, err := services.CreateNewBin()
		assert.Error(t, err)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfInsertBin: 1,
		})
	})
}

func Test_LogRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)

		db := fake.Db{
			InsertRequestFake: func(request models.Request) error {
				assert.Equal(t, id.String(), request.Bin)
				return nil
			},
		}

		services := New(&Deps{
			Db: &db,
		})

		request := models.Request{
			Bin:        id.String(),
			RecievedAt: time.Now(),
			Headers:    "headers",
			Body:       "body",
			Host:       "host",
			Method:     "method",
		}

		err = services.LogRequest(request)
		assert.NoError(t, err)

		db.VerifyCallCounts(t, &fake.Db{
			CountOfInsertRequest: 1,
		})
	})
	t.Run("invalid bin request", func(t *testing.T) {
		db := fake.Db{}
		services := New(&Deps{
			Db: &db,
		})

		request := models.Request{
			Bin:        "invalid-bin-id",
			RecievedAt: time.Now(),
			Headers:    "headers",
			Body:       "body",
			Host:       "host",
			Method:     "method",
		}

		err := services.LogRequest(request)
		assert.Error(t, err)
	})
	t.Run("error inserting request", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)

		db := fake.Db{
			InsertRequestFake: func(request models.Request) error {
				return assert.AnError
			},
		}
		services := New(&Deps{
			Db: &db,
		})

		request := models.Request{
			Bin:        id.String(),
			RecievedAt: time.Now(),
			Headers:    "headers",
			Body:       "body",
			Host:       "host",
			Method:     "method",
		}

		err = services.LogRequest(request)
		assert.Error(t, err)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfInsertRequest: 1,
		})
	})
}

func Test_GetRequestsInBin(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)

		db := fake.Db{
			GetBinContentsFake: func(binId string) ([]models.Request, error) {
				assert.Equal(t, id.String(), binId)
				return []models.Request{
					{
						Bin:        id.String(),
						RecievedAt: time.Now(),
						Headers:    "headers",
						Body:       "body",
						Host:       "host",
						Method:     "method",
					},
				}, nil
			},
		}

		services := New(&Deps{
			Db: &db,
		})

		requests, err := services.GetRequestsInBin(id.String())
		assert.NoError(t, err)
		assert.Len(t, requests, 1)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfGetBinContents: 1,
		})
	})
	t.Run("invalid bin request", func(t *testing.T) {
		db := fake.Db{}
		services := New(&Deps{
			Db: &db,
		})

		requests, err := services.GetRequestsInBin("invalid-bin-id")
		assert.Nil(t, requests)
		assert.Error(t, err)
	})
	t.Run("error getting requests", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)

		db := fake.Db{
			GetBinContentsFake: func(binId string) ([]models.Request, error) {
				assert.Equal(t, id.String(), binId)
				return nil, assert.AnError
			},
		}
		services := New(&Deps{
			Db: &db,
		})

		requests, err := services.GetRequestsInBin(id.String())
		assert.Nil(t, requests)
		assert.Error(t, err)
	})
}
