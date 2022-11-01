package main

import (
	"log"

	"user-syncer/cmd/app"
)

func main() {

	command := app.NewCommand()
	err := command.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
