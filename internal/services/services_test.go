package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	fake "app/internal/db/test"
	"app/internal/models"
)

func generateRequest() models.Request {
	headers := map[string][]string{"header": {"header"}}
	req := models.Request{
		RecievedAt: time.Now(),
		Body:       "body",
		Host:       "host",
		Method:     "method",
		Bin:        1,
	}
	_ = req.SetHeaders(headers)
	return req
}

func Test_CreateNewBin(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		generatedBinId := int64(10000)
		db := fake.Db{
			CreateBinFake: func(bin models.Bin) (int64, error) {
				return generatedBinId, nil
			},
		}

		services := New(&Deps{
			Db: &db,
		})

		binId, err := services.CreateNewBin()
		assert.NoError(t, err)
		assert.Equal(t, generatedBinId, binId)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfCreateBin: 1,
		})
	})
	t.Run("error inserting bin", func(t *testing.T) {
		db := fake.Db{
			CreateBinFake: func(bin models.Bin) (int64, error) {
				return 0, assert.AnError
			},
		}
		services := New(&Deps{
			Db: &db,
		})

		_, err := services.CreateNewBin()
		assert.Error(t, err)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfCreateBin: 1,
		})
	})
}

func Test_LogRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		request := generateRequest()

		db := fake.Db{
			InsertRequestFake: func(requestParams models.Request) error {
				assert.Equal(t, request.Bin, requestParams.Bin)
				return nil
			},
		}

		services := New(&Deps{
			Db: &db,
		})

		err := services.LogRequest(request)
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

		request := generateRequest()
		request.Bin = 0

		err := services.LogRequest(request)
		assert.Error(t, err)
	})
	t.Run("error inserting request", func(t *testing.T) {
		db := fake.Db{
			InsertRequestFake: func(request models.Request) error {
				return assert.AnError
			},
		}
		services := New(&Deps{
			Db: &db,
		})

		request := generateRequest()

		err := services.LogRequest(request)
		assert.Error(t, err)
		db.VerifyCallCounts(t, &fake.Db{
			CountOfInsertRequest: 1,
		})
	})
}

func Test_GetRequestsInBin(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		id := int64(1)

		db := fake.Db{
			GetBinContentsFake: func(binId int64) ([]models.Request, error) {
				assert.Equal(t, id, binId)
				loggedRequest := generateRequest()
				return []models.Request{loggedRequest}, nil
			},
		}

		services := New(&Deps{
			Db: &db,
		})

		requests, err := services.GetRequestsInBin(id)
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

		requests, err := services.GetRequestsInBin(0)
		assert.Nil(t, requests)
		assert.Error(t, err)
	})
	t.Run("error getting requests", func(t *testing.T) {
		id := int64(1)

		db := fake.Db{
			GetBinContentsFake: func(binId int64) ([]models.Request, error) {
				assert.Equal(t, id, binId)
				return nil, assert.AnError
			},
		}
		services := New(&Deps{
			Db: &db,
		})

		requests, err := services.GetRequestsInBin(id)
		assert.Nil(t, requests)
		assert.Error(t, err)
	})
}
