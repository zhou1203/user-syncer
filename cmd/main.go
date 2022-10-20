package main

import (
	"user-syncer/cmd/app"
)

func main() {

	command := app.NewCommand()
	err := command.Execute()
	if err != nil {
		return
	}

}
