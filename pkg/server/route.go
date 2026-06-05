package server

import (
	"net/http"
)

type Route struct {
	ServiceName string
	Path        string
	Handler     http.Handler
}
