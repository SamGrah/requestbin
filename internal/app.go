package app

import (
	"app/internal/controllers"
	"app/internal/db"
	"app/internal/router"
	"app/internal/services"
	"log"
)

type Deps struct {
	Db       Db
	Services Services
	Server   Server
}

type App struct {
	db       Db
	services Services
	server   Server
}

func NewApp() *App {
	dataService, err := db.NewDb("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	srvs := services.New(&services.Deps{
		Db: dataService,
	})

	controllers := controllers.NewControllers(&controllers.Deps{
		Services: srvs,
	})
	router := router.Routes(controllers)

	newServer := NewServer(":3000", router)

	return &App{
		db:     dataService,
		server: newServer,
	}
}

func (app *App) Init() error {
	err := app.db.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Start() error {
	err := app.server.Start()
	if err != nil {
		return err
	}

	return nil
}
