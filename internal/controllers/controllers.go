package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"app/internal/models"
	"app/internal/templates"

	"github.com/go-chi/chi/v5"
)

type Services interface {
	CreateNewBin() (int64, error)
	LogRequest(request models.Request) error
	GetRequestsInBin(binId int64) ([]models.Request, error)
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
	component := wrapComponentTemplate(
		templates.NewBin(strconv.FormatInt(binId, 10)),
		r,
	)

	err = component.Render(context.Background(), w)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/html")
}

func (c *Controllers) LogRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error reading request body: %s", err.Error())))
		return
	}

	urlBinId := chi.URLParam(r, "binId")
	binId, err := strconv.ParseInt(urlBinId, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing bin id: %s", err.Error())))
		return
	}
	reqToLog := models.Request{
		Bin:        binId,
		RecievedAt: time.Now(),
		Headers:    string(body),
		Body:       string(body),
		Host:       r.Host,
		RemoteAddr: r.RemoteAddr,
		RequestUri: r.RequestURI,
		Method:     r.Method,
	}
	reqToLog.SetHeaders(r.Header)

	err = c.services.LogRequest(reqToLog)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (c *Controllers) ViewBinContents(w http.ResponseWriter, r *http.Request) {
	log.Printf("should print view bin contents: %+v", r)
	urlBinId := chi.URLParam(r, "binId")
	binId, err := strconv.ParseInt(urlBinId, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing bin id: %s", err.Error())))
		return
	}

	requests, err := c.services.GetRequestsInBin(binId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for _, request := range requests {
		log.Printf("request: %+v", request)
	}

	reqParams := templates.ViewBinParams{
		BinId:    strconv.FormatInt(binId, 10),
		Hostname: r.Host,
		Requests: requests,
	}
	component := templates.Layout(templates.ViewBinContents(reqParams))
	log.Printf("should print view bin html: %+v", component)

	err = component.Render(context.Background(), w)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/html")
}
