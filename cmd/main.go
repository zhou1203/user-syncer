package main

import (
	"user-generator/cmd/app"
)

func main() {

	command := app.NewCommand()
	err := command.Execute()
	if err != nil {
		return
	}

}
