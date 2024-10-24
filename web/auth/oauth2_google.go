package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

type OAuth2Google struct {
}

func (oauth2google OAuth2Google) StartFlow() (*http.Client, error) {
	config, err := getConfig("credentials.json")
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to parse client secret file to config: %v", err)
		return nil, errors.New(errorMessage)
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	tokFile := "token.json"
	tok, err := getTokenFromFile(tokFile)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		log.Println("Opening the browser and start the OAuth2 flow.")

		if openBrowserErr := openBrowser(authURL); openBrowserErr != nil {
			errorMessage := fmt.Sprintf("OAuth2 flow was disrupted by: %v", openBrowserErr)
			return nil, errors.New(errorMessage)
		}

		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

func (oauth2google OAuth2Google) StoreToken(authCode string, path string) error {
	config, err := getConfig("credentials.json")
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to parse client secret file to config: %v", err)
		return errors.New(errorMessage)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		errorMessage := fmt.Sprintf("Error converting an authorization code into a token %v", err)
		return errors.New(errorMessage)
	}

	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to cache oauth token: %v", err)
		return errors.New(errorMessage)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(tok)
	return nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	cmd = "rundll32"
	args = append(args, "url.dll,FileProtocolHandler", url)

	return exec.Command(cmd, args...).Start()
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getConfig(credentialsPath string) (*oauth2.Config, error) {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to read client secret file: %v", err)
		return nil, errors.New(errorMessage)
	}

	// If modifying these scopes, delete your previously saved token.json.
	return google.ConfigFromJSON(b, sheets.SpreadsheetsScope, drive.DriveScope)
}
