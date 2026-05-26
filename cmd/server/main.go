package main

import (
	"context"
	"fmt"
	"net/http"

	internalfx "e2b/internal"
	pkgfx "e2b/pkg"
	"e2b/pkg/config"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		internalfx.Module,
		pkgfx.Module,
		fx.Invoke(StartHTTPServer),
	)
	app.Run()
}

type StartHTTPServerParams struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    config.Config
	Mux       *http.ServeMux
}

func StartHTTPServer(in StartHTTPServerParams) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", in.Config.Port),
		Handler: in.Mux,
	}

	in.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		},
	})
	return nil
}
