package main

import (
	"fmt"
	"os"

	"github.com/usace/nsi-cli/commands"
	"github.com/usace/nsi-cli/config"
)

func main() {
	// If no sub-command is provided, print a usage message and exit.
	if len(os.Args) == 1 {
		println("Usage: nsi-cli <commands> [<args>]")
		return
	}

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error: Could not load config file")
		os.Exit(3)
	}

	var cmd commands.Command

	// Determine which sub-command should be run, and parse the remaining arguments.
	switch os.Args[1] {
	case "auth":
		cmd = commands.NewAuthCommand(config)

	case "groups":
		cmd = commands.NewGroupsCommand(config)

	case "users":
		cmd = commands.NewUsersCommand(config)

	default:
		fmt.Printf("%q is not a valid command\n", os.Args[1])
		os.Exit(2)
	}

	// Parse the rest of the command-line arguments.
	err = cmd.Parse(os.Args[2:])
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
