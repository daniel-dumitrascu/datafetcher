package main

import (
	"datafetcher/extern/cloud"
	"datafetcher/jobs"
	"datafetcher/web"
	"datafetcher/web/auth"
	"errors"
	"fmt"
	"net/http"
)

func Init(path string) error {
	web.RegisterHandlers()

	oauth := auth.CreateOAuth2(cloud.Google)
	err := oauth.StartFlow()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to authorize the use of google api: %v", err)
		return errors.New(errorMessage)
	}

	//TODO New method that gets a client when needed
	//oauth.GetClient()

	jobs.StartJobs(path)
	return nil
}

func Run() {
	http.ListenAndServe(":6788", nil)
}
