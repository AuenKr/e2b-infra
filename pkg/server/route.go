package server

import (
	"net/http"
)

type Route struct {
	// ServiceName from generated protofile for registering in reflection
	ServiceName string
	// Service handler path and handler
	Path    string
	Handler http.Handler
}
