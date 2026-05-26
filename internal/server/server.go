package server

import (
	"net/http"

	"e2b/pkg/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RegisterRoutesParams struct {
	fx.In
	Routers []server.Route `group:"grpc-routes"`
	Logger  *zap.Logger
}

func NewServerMux(in RegisterRoutesParams) *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range in.Routers {
		in.Logger.Info("Registering route", zap.String("path", r.Path))
		mux.Handle(r.Path, r.Handler)
	}
	return mux
}
