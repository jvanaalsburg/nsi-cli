package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"syscall"

	"github.com/usace/nsi-cli/api"
	"github.com/usace/nsi-cli/config"
	"golang.org/x/term"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthCommand struct {
	config        config.Config
	loginFlags    *flag.FlagSet
	statusFlags   *flag.FlagSet
	tokenFlags    *flag.FlagSet
	loginEmail    *string
	loginPassword *string
}

func NewAuthCommand(config config.Config) AuthCommand {
	loginFlags := flag.NewFlagSet("login", flag.ExitOnError)
	statusFlags := flag.NewFlagSet("status", flag.ExitOnError)
	tokenFlags := flag.NewFlagSet("token", flag.ExitOnError)

	return AuthCommand{
		config:        config,
		loginFlags:    loginFlags,
		statusFlags:   statusFlags,
		tokenFlags:    tokenFlags,
		loginEmail:    loginFlags.String("email", "", "Account email address"),
		loginPassword: loginFlags.String("password", "", "Account password"),
	}
}

func (c AuthCommand) Parse(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You must specify an action.")
	}

	action := args[0]
	switch action {
	case "login":
		c.loginFlags.Parse(args[1:])

	case "status":
		c.statusFlags.Parse(args[1:])

	case "token":
		c.tokenFlags.Parse(args[1:])

	default:
		return fmt.Errorf("Invalid auth action: %s", action)
	}

	return nil
}

func (c AuthCommand) Validate() error {
	if c.loginFlags.Parsed() {
		if *c.loginEmail == "" {
			return fmt.Errorf("You must specify an email address")
		}
	}

	return nil
}

func (c AuthCommand) Exec() {
	if c.loginFlags.Parsed() {
		c.login()
	}

	if c.statusFlags.Parsed() {
		c.authStatus()
	}

	if c.tokenFlags.Parsed() {
		c.authToken()
	}

	return
}

func (c AuthCommand) login() {
	log.Println("logging in...")

	var password string = *c.loginPassword

	if password == "" {
		var err error
		password, err = getPassword()
		if err != nil {
			log.Fatalf("Error getting credentials: %v", err)
		}
	}

	// Attempt to login with the provided email and password.
	token, err := authenticate(c.config, *c.loginEmail, password)
	if err != nil {
		log.Fatalf("Error requesting auth token: %v", err)
	}

	// Update the auth data and save the config file.
	c.config.Auth.Email = *c.loginEmail
	c.config.Auth.Token = token

	err = c.config.SaveConfig()
	if err != nil {
		log.Fatalf("Error saving config file: %v", err)
	}
}

func authenticate(config config.Config, username, password string) (string, error) {
	data := url.Values{
		"username": {username},
		"password": {password},
	}

	res, err := api.PostForm(config, data, "login")
	if err != nil {
		return "", err
	}

	var response AuthResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func getPassword() (string, error) {
	fmt.Print("Enter password....: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	password := string(bytes)

	return password, nil
}

func (c AuthCommand) authStatus() {
	if c.config.Auth.Token == "" {
		fmt.Printf("Not currently logged in.")
	} else {
		fmt.Printf("Logged into account %s\n", c.config.Auth.Email)
	}
}

func (c AuthCommand) authToken() {
	if c.config.Auth.Token == "" {
		log.Fatalf("No auth token found")
	}

	fmt.Print(c.config.Auth.Token)
}
