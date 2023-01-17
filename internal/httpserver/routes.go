package httpserver

import "net/http"

type Route struct {
	Verb string
	Path string
	Fn   func(http.ResponseWriter, *http.Request)
}

type Routes []Route
