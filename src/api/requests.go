package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/usace/nsi-cli/config"
)

func Get(config config.Config, path ...string) (*http.Response, error) {
	// Construct the API endpoint from the path arguments.
	endpoint := strings.Join(path, "/")

	// Initialize a new GET request.
	url := fmt.Sprintf("%s/%s", config.Api.UrlRoot, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add the authorization header.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Auth.Token))

	// Make the API request.
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func PostJson(config config.Config, payload []byte, path ...string) (*http.Response, error) {
	// Construct the API endpoint from the path arguments.
	endpoint := strings.Join(path, "/")

	url := fmt.Sprintf("%s/%s", config.Api.UrlRoot, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Auth.Token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
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

func Delete(config config.Config, path ...string) (*http.Response, error) {
	// Construct the API endpoint from the path arguments.
	endpoint := strings.Join(path, "/")

	url := fmt.Sprintf("%s/%s", config.Api.UrlRoot, endpoint)
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Auth.Token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
