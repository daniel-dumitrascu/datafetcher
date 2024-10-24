package main

import (
	"datafetcher/jobs"
	"datafetcher/web"
	"net/http"
)

func Init(path string) error {
	web.RegisterHandlers()
	jobs.StartJobs(path)
	return nil
}

func Run() {
	http.ListenAndServe(":6788", nil)
}
