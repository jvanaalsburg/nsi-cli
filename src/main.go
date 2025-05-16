package main

import (
	"fmt"
	"os"

	"github.com/usace/nsi-cli/commands"
)

func main() {
	// If no sub-command is provided, print a usage message and exit.
	if len(os.Args) == 1 {
		println("Usage: nsi-cli <commands> [<args>]")
		return
	}

	var cmd commands.Command

	// Determine which sub-command should be run, and parse the remaining arguments.
	switch os.Args[1] {
	case "auth":
		cmd = commands.NewAuthCommand()

	case "users":
		cmd = commands.NewUsersCommand()

	default:
		fmt.Printf("%q is not a valid command\n", os.Args[1])
		os.Exit(2)
	}

	// Parse the rest of the command-line arguments.
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		println(err.Error())
		return
	}

	// Validate the command-line arguments.
	err = cmd.Validate()
	if err != nil {
		println(err.Error())
		return
	}

	// Run the command.
	cmd.Exec()
}
