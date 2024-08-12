package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handlers interface {
	Index(w http.ResponseWriter, r *http.Request)
	NewBin(w http.ResponseWriter, r *http.Request)
	Intro(w http.ResponseWriter, r *http.Request)
}

func Routes(h Handlers) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	router.Group(func(router chi.Router) {
		router.Get("/", h.Index)
		router.Get("/intro", h.Intro)
		router.Get("/new-bin", h.NewBin)
		// router.HandleFunc("/bin/{binId}", h.BinMgmtHandler.LogRequest)
		// router.Get("/bin-contents/{binId}", h.BinMgmtHandler.FetchBinContents)
	})

	return router
}
