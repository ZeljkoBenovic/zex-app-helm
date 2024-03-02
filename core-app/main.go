package main

import (
	"log"

	"be/pkg/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Fatalln(a.Run())
}
