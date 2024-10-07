package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// DownloadFile downloads a file from Google Drive
func DownloadFile(driveService *drive.Service, fileId string, destination string) error {
	resp, err := driveService.Files.Export(fileId, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet").Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

func UpdateFile(driveService *sheets.Service, fileId string) error {
	// The range to update (e.g., "Sheet1!A1")
	rangeToUpdate := "Sheet1!A1"

	// Value to write into the cell
	value := [][]interface{}{
		{"Dumitrascu"}, // This will be written into cell A1
	}

	// Prepare the update request
	valueRange := &sheets.ValueRange{
		Values: value,
	}

	// Call the Update method
	_, err := driveService.Spreadsheets.Values.Update(fileId, rangeToUpdate, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	fileID := "15Yt9IWfuiASjvv38btbelcoglLx6RnC-cXuds7lQLQg"

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	file, err := srv.Files.Get(fileID).Fields("id, name, mimeType, size").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve file: %v", err)
	}

	// Output the file metadata
	fmt.Printf("File ID: %s\n", file.Id)
	fmt.Printf("File Name: %s\n", file.Name)
	fmt.Printf("MIME Type: %s\n", file.MimeType)
	fmt.Printf("File Size: %d bytes\n", file.Size)

	// localFilePath := "C:\\Users\\DanielDumitrascu\\Desktop\\algo\\1\\excels\\example.xlsx"

	// // Download the Excel file from Google Drive
	// if err := DownloadFile(srv, fileID, localFilePath); err != nil {
	// 	log.Fatalf("Unable to download file: %v", err)
	// }
	// fmt.Println("File downloaded successfully.")

	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	if err := UpdateFile(sheetsSrv, fileID); err != nil {
		log.Fatalf("Unable to download file: %v", err)
	}

	fmt.Println("File was updated successfully.")
}
