package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/usace/nsi-cli/config"
)

func Get(config config.Config, path ...string) (string, error) {
	// Construct the API endpoint from the path arguments.
	endpoint := strings.Join(path, "/")

	// Initialize a new GET request.
	url := fmt.Sprintf("%s/%s", config.Api.UrlRoot, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// Add the authorization header.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Auth.Token))

	// Make the API request.
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Return the response body as a string.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PostForm(config config.Config, data url.Values, path ...string) (*http.Response, error) {
	// Construct the API endpoint from the path arguments.
	endpoint := strings.Join(path, "/")

	url := fmt.Sprintf("%s/%s", config.Api.UrlRoot, endpoint)
	res, err := http.PostForm(url, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
