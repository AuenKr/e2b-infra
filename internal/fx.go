// Package internalfx: Internal module for the server.
package internalfx

import (
	"e2b/internal/sandbox"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"server",
	fx.Provide(sandbox.NewSandboxServer),
)
