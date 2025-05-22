package commands

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/usace/nsi-cli/api"
	"github.com/usace/nsi-cli/config"
	"github.com/usace/nsi-cli/utils"
)

type UsersCommand struct {
	config              config.Config
	listUsersFlags      *flag.FlagSet
	findUsersFlags      *flag.FlagSet
	findUsersId         *string
	createUserFlags     *flag.FlagSet
	createUserEmail     *string
	createUserFirstName *string
	createUserLastName  *string
	createUserPassword  *string
}

func NewUsersCommand(config config.Config) UsersCommand {
	listUsersFlags := flag.NewFlagSet("list", flag.ExitOnError)
	findUsersFlags := flag.NewFlagSet("find", flag.ExitOnError)
	createUserFlags := flag.NewFlagSet("create", flag.ExitOnError)

	return UsersCommand{
		config:              config,
		listUsersFlags:      listUsersFlags,
		findUsersFlags:      findUsersFlags,
		findUsersId:         findUsersFlags.String("user-id", "", "The user ID"),
		createUserFlags:     createUserFlags,
		createUserEmail:     createUserFlags.String("email", "", "Email address"),
		createUserFirstName: createUserFlags.String("first-name", "", "First Name"),
		createUserLastName:  createUserFlags.String("last-name", "", "Last Name"),
		createUserPassword:  createUserFlags.String("password", utils.RandomPassword(), "Password"),
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

	case "create":
		c.createUserFlags.Parse(args[1:])

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

	if c.createUserFlags.Parsed() {
		c.createUser()
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

func (c UsersCommand) createUser() {
	data := url.Values{
		"username":   {*c.createUserEmail},
		"first_name": {*c.createUserFirstName},
		"last_name":  {*c.createUserLastName},
		"password":   {*c.createUserPassword},
	}

	log.Printf("Creating account %s (password: %s)",
		*c.createUserEmail,
		*c.createUserPassword,
	)

	res, err := api.PostForm(c.config, data, "register")
	if err != nil {
		log.Fatalf("Error registering user: %v", err)
	}

	fmt.Printf(api.ResponseStr(res))
}
