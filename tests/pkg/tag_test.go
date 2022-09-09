package pkg_test

import (
	"testing"

	h "github.com/kaankoken/helper/pkg"
	"github.com/kaankoken/versioning-tool/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
)

func TestTag(t *testing.T) {
	t.Run("tag=invoke-method", func(t *testing.T) {
		tag := pkg.ConstructTag()

		assert.NotNil(t, tag)
		assert.NotNil(t, &tag)
	})
}

func TestTagWithFx(t *testing.T) {
	t.Run("tag=injection-test", func(t *testing.T) {
		var g fx.DotGraph

		app := fxtest.New(
			t,
			fx.Logger(fxtest.NewTestPrinter(t)),
			fx.WithLogger(func() fxevent.Logger { return fxtest.NewTestLogger(t) }),
			fx.Populate(&g),
			pkg.TagModule,
		).RequireStart()

		defer app.RequireStop()

		require.NoError(t, app.Err())
		assert.Contains(t, g, `"fx.DotGraph" [label=<fx.DotGraph>];`)
	})

	t.Run("tag=injection-test-with-functions", func(t *testing.T) {
		var g fx.DotGraph
		var l *h.Tag

		app := fxtest.New(
			t,
			fx.Logger(fxtest.NewTestPrinter(t)),
			fx.WithLogger(func() fxevent.Logger { return fxtest.NewTestLogger(t) }),
			fx.Populate(&g),
			fx.Populate(&l),
			pkg.TagModule,
		).RequireStart()

		defer app.RequireStop()

		require.NoError(t, app.Err())
		assert.Contains(t, g, `"fx.DotGraph" [label=<fx.DotGraph>];`)
	})
}
