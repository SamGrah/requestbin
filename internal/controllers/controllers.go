package controllers

import (
	"context"
	"log"
	"net/http"

	"app/internal/templates"
)

type Controllers struct{}

type Deps struct{}

func NewControllers(deps *Deps) *Controllers {
	return &Controllers{}
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
	component := wrapComponentTemplate(templates.NewBin("12345"), r)

	err := component.Render(context.Background(), w)
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
