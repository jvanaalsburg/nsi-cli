package api

import (
	"io"
	"net/http"
)

func ResponseStr(res *http.Response) string {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	return string(body)
}
