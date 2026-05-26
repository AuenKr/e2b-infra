package server

import (
	"net/http"

	"e2b/pkg/config"
	"e2b/pkg/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHTTPMux(config config.Config) *http.ServeMux {
	return http.NewServeMux()
}

type RegisterRoutesParams struct {
	fx.In
	Mux     *http.ServeMux
	Routers []server.Route `group:"grpc-routes"`
	Logger  *zap.Logger
}

func RegisterRoutes(in RegisterRoutesParams) *http.ServeMux {
	for _, r := range in.Routers {
		in.Logger.Info("Registering route", zap.String("path", r.Path))
		in.Mux.Handle(r.Path, r.Handler)
	}
	return in.Mux
}
