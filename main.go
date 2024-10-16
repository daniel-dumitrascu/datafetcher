package main

import (
	"net/http"

	"test/web"
)

func main() {
	web.RegisterHandlers()
	http.ListenAndServe(":6788", nil)
}
