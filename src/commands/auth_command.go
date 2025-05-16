package commands

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthCommand struct {
	flags *flag.FlagSet
}

func NewAuthCommand() AuthCommand {
	flags := flag.NewFlagSet("login", flag.ExitOnError)

	return AuthCommand{
		flags: flags,
	}
}

func (c AuthCommand) Parse(args []string) {
	c.flags.Parse(args)
}

func (c AuthCommand) Validate() error {
	actions := []string{"login", "status"}

	if len(c.flags.Args()) == 0 {
		return fmt.Errorf("You must specify an action.")
	}

	action := c.flags.Arg(0)
	if !slices.Contains(actions, action) {
		return fmt.Errorf("Invalid action")
	}

	return nil
}

func (c AuthCommand) Exec() {
	action := c.flags.Arg(0)

	switch action {
	case "login":
		login()

	case "status":
		authStatus()
	}

	return
}

func login() {
	log.Println("logging in...")

	email, password, err := getCredentials()
	if err != nil {
		log.Fatalf("Error getting credentials: %v", err)
	}

	token, err := requestToken(email, password)
	if err != nil {
		log.Fatalf("Error requesting auth token: %v", err)
	}

	fmt.Printf("token: %s", token)
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

func getCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email.......: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter password....: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", "", err
	}

	email = strings.TrimRight(email, "\n")
	password := string(bytes)

	return email, password, nil
}

func authStatus() {
	log.Println("checking auth status...")
}
