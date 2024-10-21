package main

import (
	"datafetcher/web"
	"net/http"
)

func Init() error {
	web.RegisterHandlers()
	return nil
}

func Run() {
	http.ListenAndServe(":6788", nil)
}
