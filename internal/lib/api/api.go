package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusOk   = errors.New("invalid status code")
	ErrInvalidStatusCode = errors.New("invalis status code")
)

func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrInvalidStatusCode
	}

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidStatusOk)
	}

	return resp.Header.Get("Location"), nil
}
