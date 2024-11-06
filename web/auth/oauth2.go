package auth

import "net/http"

type OAuth2 interface {
	StartFlow() error
	StoreToken(authCode string, path string) error
	GetClient() *http.Client
}
