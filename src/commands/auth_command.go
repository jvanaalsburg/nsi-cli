package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"syscall"

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
	loginEmail    *string
	loginPassword *string
}

func NewAuthCommand() AuthCommand {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Errorf("Error loading config file: %v", err)
	}

	loginFlags := flag.NewFlagSet("login", flag.ExitOnError)
	statusFlags := flag.NewFlagSet("status", flag.ExitOnError)

	return AuthCommand{
		config:        config,
		loginFlags:    loginFlags,
		statusFlags:   statusFlags,
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

	token, err := c.requestToken(password)
	if err != nil {
		log.Fatalf("Error requesting auth token: %v", err)
	}

	err = c.saveToken(token)
	if err != nil {
		log.Fatalf("Error saving config file: %v", err)
	}
}

func (c AuthCommand) saveToken(token string) error {
	c.config.Auth.Email = *c.loginEmail
	c.config.Auth.Token = token

	return c.config.Save()
}

func (c AuthCommand) requestToken(password string) (string, error) {
	data := url.Values{
		"username": {*c.loginEmail},
		"password": {password},
	}

	url := fmt.Sprintf("%s/login", c.config.Api.UrlRoot)
	req, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}

	var response AuthResponse
	err = json.NewDecoder(req.Body).Decode(&response)
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
	log.Println("checking auth status...")

	if c.config.Auth.Token == "" {
		fmt.Printf("Not currently logged in.")
	} else {
		fmt.Printf("Logged into account %s\n", c.config.Auth.Email)
	}
}
