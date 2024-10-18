package auth

import "net/http"

type OAuth2 interface {
	StartFlow() (*http.Client, error)
	StoreToken(authCode string, path string) error
}
