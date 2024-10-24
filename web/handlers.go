package web

import (
	"context"
	"datafetcher/extern/cloud"
	"datafetcher/web/auth"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func testUpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure that only POST requests are handled
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Use OAuth2 to login to google in order to use his APIs
	oauth := auth.CreateOAuth2(cloud.Google)
	client, err := oauth.StartFlow()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to authorize the use of google api: %v", err)
		http.Error(w, errorMessage, http.StatusUnauthorized)
		return
	}

	fileID := "15Yt9IWfuiASjvv38btbelcoglLx6RnC-cXuds7lQLQg"

	ctx := context.Background()
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
	// if err := downloadFile(srv, fileID, localFilePath); err != nil {
	// 	log.Fatalf("Unable to download file: %v", err)
	// }
	// fmt.Println("File downloaded successfully.")

	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	if err := updateFile(sheetsSrv, fileID); err != nil {
		log.Fatalf("Unable to download file: %v", err)
	}

	fmt.Println("File was updated successfully.")

	// Parse the incoming JSON
	// var data PostData
	// err := json.NewDecoder(r.Body).Decode(&data)
	// if err != nil {
	// 	http.Error(w, "Error parsing request body", http.StatusBadRequest)
	// 	return
	// }

	// Respond with a confirmation message
	//response := fmt.Sprintf("Received: Name = %s, Value = %d", data.Name, data.Value)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Endpoint was succesfully called"))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters from the URL
	queryParams := r.URL.Query()

	code := queryParams.Get("code")
	if code == "" {
		http.Error(w, "Missing 'code' query parameter", http.StatusBadRequest)
		return
	}

	// Store it
	path := "token.json"
	oauth := auth.CreateOAuth2(cloud.Google)
	oauth.StoreToken(code, path)

	fmt.Printf("Token has been stored at location: %v\n", path)
}

func RegisterHandlers() {
	http.HandleFunc("/api/v1/testupdatedata", testUpdateDataHandler)
	http.HandleFunc("/api/v1/oauth2/token", tokenHandler)
}

// downloadFile downloads a file from Google Drive
func downloadFile(driveService *drive.Service, fileId string, destination string) error {
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

func updateFile(driveService *sheets.Service, fileId string) error {
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
