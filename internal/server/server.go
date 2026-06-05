package server

import (
	"net/http"

	"e2b/pkg/config"
	"e2b/pkg/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RegisterRoutesParams struct {
	fx.In
	Config           config.Config
	Routers          []server.Route `group:"grpc-routes"`
	ReflectionRoutes []server.Route `group:"reflection"`
	Logger           *zap.Logger
}

func NewServerMux(in RegisterRoutesParams) *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range in.Routers {
		in.Logger.Info("Registering route", zap.String("path", r.Path))
		mux.Handle(r.Path, r.Handler)
	}

	if in.Config.Mode == "dev" {
		for _, r := range in.ReflectionRoutes {
			in.Logger.Info("Registering reflection route", zap.String("path", r.Path))
			mux.Handle(r.Path, r.Handler)
		}
	}
	return mux
}
