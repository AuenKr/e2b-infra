package main

import (
	"context"
	"fmt"

	internalfx "e2b/internal"

	"go.uber.org/fx"
)

func main() {
	// Server will start here
	fmt.Println("Server is running")

	app := fx.New(
		fx.Provide(internalfx.Module),
		fx.Invoke(),
	)
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}
