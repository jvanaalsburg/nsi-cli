package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/usace/nsi-cli/api"
	"github.com/usace/nsi-cli/config"
)

type UsersCommand struct {
	config         config.Config
	listUsersFlags *flag.FlagSet
	findUsersFlags *flag.FlagSet
	findUsersId    *string
}

func NewUsersCommand(config config.Config) UsersCommand {
	listUsersFlags := flag.NewFlagSet("list", flag.ExitOnError)
	findUsersFlags := flag.NewFlagSet("find", flag.ExitOnError)

	return UsersCommand{
		config:         config,
		listUsersFlags: listUsersFlags,
		findUsersFlags: findUsersFlags,
		findUsersId:    findUsersFlags.String("user-id", "", "The user ID"),
	}
}

func (c UsersCommand) Parse(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You must specify an action.")
	}

	action := args[0]
	switch action {
	case "list":
		c.listUsersFlags.Parse(args[1:])

	case "find":
		c.findUsersFlags.Parse(args[1:])

	default:
		return fmt.Errorf("Invalid users action: %s", action)
	}

	return nil
}

func (c UsersCommand) Validate() error {
	return nil
}

func (c UsersCommand) Exec() {
	if c.listUsersFlags.Parsed() {
		c.findAllUsers()
	}

	if c.findUsersFlags.Parsed() {
		c.findUser()
	}

	return
}

func (c UsersCommand) findAllUsers() {
	users, err := api.Get(c.config, "users")
	if err != nil {
		log.Fatalf("Error fetching user records: %v", err)
	}

	fmt.Printf(api.ResponseStr(users))
}

func (c UsersCommand) findUser() {
	user, err := api.Get(c.config, "users", *c.findUsersId)
	if err != nil {
		log.Fatalf("Error fetching user record: %v", err)
	}

	fmt.Printf(api.ResponseStr(user))
}
