package pkg

import (
	"github.com/kaankoken/helper/pkg"
	"go.uber.org/fx"
)

// FlavorModule -> Dependency Injection for FlavorModule module
var FlavorModule = fx.Options(fx.Provide(ConstructFlavor))

// ConstructFlavor -> Flavor needed by DI & {helper.Logger}
func ConstructFlavor() *pkg.Flavor {
	return &pkg.Flavor{F: "dev"}
}
