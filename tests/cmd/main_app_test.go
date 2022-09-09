package cmd_test

import (
	"context"
	"testing"

	"github.com/kaankoken/helper/pkg/helper"
	"github.com/kaankoken/versioning-tool/cmd"
	"github.com/kaankoken/versioning-tool/pkg"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestMainApp(t *testing.T) {
	t.Parallel()

	t.Run("test-main-app-func", func(t *testing.T) {
		t.Parallel()

		app := fxtest.New(
			t,
			pkg.FlavorModule,
			pkg.TagModule,
			helper.LoggerModule,
			fx.Invoke(cmd.RegisterHooks),
		)

		app.RequireStart()
		defer app.RequireStop()
	})

	t.Run("test-main-app=no-args", func(t *testing.T) {
		t.Parallel()

		_, err := cmd.MainApp(make([]string, 0))

		assert.NotNil(t, err)
	})

	t.Run("test-main-app=one-args", func(t *testing.T) {
		t.Parallel()

		_, err := cmd.MainApp([]string{""})

		assert.NotNil(t, err)
	})

	t.Run("test-main-app=without-base", func(t *testing.T) {
		t.Parallel()

		app, err := cmd.MainApp([]string{"key", "owner", "repo"})

		assert.Nil(t, err)
		app.Start(context.Background())

		defer app.Stop(context.Background())
	})

	t.Run("test-main-app=with-base", func(t *testing.T) {
		t.Parallel()

		app, err := cmd.MainApp([]string{"key", "owner", "repo", "master"})

		assert.Nil(t, err)
		app.Start(context.Background())

		defer app.Stop(context.Background())
	})
}
