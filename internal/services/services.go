package services

import (
	"app/internal/models"

	"github.com/google/uuid"
)

type Db interface {
	InsertBin(bin models.Bin) error
	InsertRequest(request models.Request) error
	GetBinContents(binId string) ([]models.Request, error)
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

func (s *Services) CreateNewBin() (string, error) {
	binId, err := uuid.NewV7() // uses ms since epoch, so it's unique :arm:
	if err != nil {
		return "", err
	}

	bin := models.Bin{
		BinId: binId.String(),
	}

	err = s.db.InsertBin(bin)
	if err != nil {
		return "", err
	}

	return binId.String(), nil
}

func (s *Services) LogRequest(request models.Request) error {
	if err := BinIdValidation(request.Bin); err != nil {
		return err
	}

	return s.db.InsertRequest(request)
}

func (s *Services) GetRequestsInBin(binId string) ([]models.Request, error) {
	if err := BinIdValidation(binId); err != nil {
		return nil, err
	}
	return s.db.GetBinContents(binId)
}

func BinIdValidation(binId string) error {
	_, err := uuid.Parse(binId)
	if err != nil {
		return err
	}

	return nil
}
