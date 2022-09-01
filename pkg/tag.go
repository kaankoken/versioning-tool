package pkg

import (
	"github.com/kaankoken/helper/pkg"
	"go.uber.org/fx"
)

// TagModule -> Dependency Injection for TagModule module
var TagModule = fx.Options(fx.Provide(ConstructTag))

// ConstructTag -> Creating Tag needed by DI & {helper.Logger}
func ConstructTag() *pkg.Tag {
	tag := "Versioning Tool -> "

	return &pkg.Tag{T: tag}
}
