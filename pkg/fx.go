// Package pkgfx: Package module for the pkg.
package pkgfx

import (
	"e2b/pkg/config"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"package",
	fx.Provide(config.NewConfig),
	fx.Provide(config.NewLogger),
	fx.Provide(config.NewK8sClusterClient),
)
