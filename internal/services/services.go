package services

import (
	"app/internal/models"
	"fmt"
)

type Db interface {
	CreateBin(bin models.Bin) (int64, error)
	InsertRequest(request models.Request) error
	GetBinContents(binId int64) ([]models.Request, error)
}

type Services struct {
	db Db
}

type Deps struct {
	Db Db
}

func New(deps *Deps) *Services {
	return &Services{
		db: deps.Db,
	}
}

func (s *Services) CreateNewBin() (int64, error) {
	binId, err := s.db.CreateBin(models.Bin{})
	if err != nil {
		return 0, err
	}

	return binId, nil
}

func (s *Services) LogRequest(request models.Request) error {
	if err := BinIdValidation(request.Bin); err != nil {
		return err
	}

	return s.db.InsertRequest(request)
}

func (s *Services) GetRequestsInBin(binId int64) ([]models.Request, error) {
	if err := BinIdValidation(binId); err != nil {
		return nil, err
	}
	return s.db.GetBinContents(binId)
}

func BinIdValidation(binId int64) error {
	if binId <= 0 {
		return fmt.Errorf("invalid bin id: %d", binId)
	}

	return nil
}
