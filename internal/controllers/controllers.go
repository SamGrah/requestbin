package controllers

import (
	"context"
	"log"
	"net/http"

	"app/internal/models"
	"app/internal/templates"
)

type Services interface {
	CreateNewBin() (string, error)
	LogRequest(request models.Request) error
	GetRequestsInBin(binId string) ([]models.Request, error)
}

type Controllers struct {
	services Services
}

type Deps struct {
	Services Services
}

func NewControllers(deps *Deps) *Controllers {
	return &Controllers{
		services: deps.Services,
	}
}

func (c *Controllers) Index(w http.ResponseWriter, r *http.Request) {
	component := templates.Layout(templates.Intro())

	err := component.Render(context.Background(), w)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/html")
}

func (c *Controllers) NewBin(w http.ResponseWriter, r *http.Request) {
	binId, err := c.services.CreateNewBin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	component := wrapComponentTemplate(templates.NewBin(binId), r)

	err = component.Render(context.Background(), w)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/html")
}

func (c *Controllers) Intro(w http.ResponseWriter, r *http.Request) {
	component := wrapComponentTemplate(templates.Intro(), r)

	err := component.Render(context.Background(), w)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/html")
}

func (c *Controllers) ViewBinContents(w http.ResponseWriter, r *http.Request) {}
