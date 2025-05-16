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
	loginFlags    *flag.FlagSet
	statusFlags   *flag.FlagSet
	loginEmail    *string
	loginPassword *string
}

func NewAuthCommand() AuthCommand {
	loginFlags := flag.NewFlagSet("login", flag.ExitOnError)
	statusFlags := flag.NewFlagSet("status", flag.ExitOnError)

	return AuthCommand{
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
		authStatus()
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

	token, err := requestToken(*c.loginEmail, password)
	if err != nil {
		log.Fatalf("Error requesting auth token: %v", err)
	}

	err = saveToken(*c.loginEmail, token)
	if err != nil {
		log.Fatalf("Error saving config file: %v", err)
	}
}

func saveToken(email, token string) error {
	var config config.Config
	config.Auth.Email = email
	config.Auth.Token = token

	return config.Save()
}

func requestToken(username, password string) (string, error) {
	data := url.Values{
		"username": {username},
		"password": {password},
	}

	url := "http://api:4141/nsiapi/login"
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

func authStatus() {
	log.Println("checking auth status...")
}
