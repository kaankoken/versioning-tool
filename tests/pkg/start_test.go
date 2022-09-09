package pkg_test

import (
	"testing"

	"github.com/kaankoken/versioning-tool/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
)

func TestStart(t *testing.T) {
	t.Parallel()

	var structCompareType pkg.InputStruct

	t.Run("start=without-base", func(t *testing.T) {
		t.Parallel()

		input := pkg.Run("key", "owner", "repo")

		assert.NotNil(t, input)
		assert.IsType(t, &structCompareType, &input)
	})

	t.Run("start=empty-string-base", func(t *testing.T) {
		t.Parallel()

		input := pkg.Run("key", "owner", "repo", "")

		assert.NotNil(t, input)
		assert.IsType(t, &structCompareType, &input)
	})

	t.Run("start=with-base", func(t *testing.T) {
		t.Parallel()

		input := pkg.Run("key", "owner", "repo", "master")

		assert.NotNil(t, input)
		assert.IsType(t, &structCompareType, &input)
	})
}

func TestStartWithFx(t *testing.T) {
	t.Run("start=injection-test", func(t *testing.T) {
		var g fx.DotGraph
		var l *pkg.InputStruct

		app := fxtest.New(
			t,
			fx.Logger(fxtest.NewTestPrinter(t)),
			fx.WithLogger(func() fxevent.Logger { return fxtest.NewTestLogger(t) }),
			fx.Populate(&g),
			fx.Supply(pkg.Run("key", "owner", "repo", "master")),
			fx.Populate(&l),
		).RequireStart()

		defer app.RequireStop()

		require.NoError(t, app.Err())
		assert.Contains(t, g, `"fx.DotGraph" [label=<fx.DotGraph>];`)
	})
}
