package main

import (
	"user-export/cmd"
)

func main() {

	command := cmd.NewCommand()
	err := command.Execute()
	if err != nil {
		return
	}

}
