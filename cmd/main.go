package main

import (
	"log"

	app "app/internal"
)

func main() {
	myApp := app.NewApp()
	err := myApp.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = myApp.Start()
	if err != nil {
		log.Fatal(err)
	}
}
