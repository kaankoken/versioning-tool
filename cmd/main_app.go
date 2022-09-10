package cmd

import (
	"context"
	"fmt"

	"github.com/kaankoken/helper/pkg/helper"
	"github.com/kaankoken/versioning-tool/pkg"
	versionlabel "github.com/kaankoken/versioning-tool/pkg/version-label"
	"go.uber.org/fx"
)

// MainApp -> Registering main app to FX
func MainApp(args []string) (app *fx.App, err error) {
	// argument controls
	if len(args) <= 0 {
		return nil, fmt.Errorf("%s", "Wrong number of inputs, no arguments found")
	}

	if len(args) < 3 {
		return nil, fmt.Errorf("%s", "Wrong number of inputs, cannot be less than 3")
	}

	input := pkg.Run(args[0], args[1], args[2], args[3:]...)

	a := fx.New(
		fx.Supply(input),
		pkg.TagModule,
		pkg.FlavorModule,
		helper.LoggerModule,
		versionlabel.GithubClient,
		versionlabel.VersionLabelModule,
		fx.Invoke(RegisterHooks),
	)

	return a, nil
}

// RegisterHooks -> Registering lifecycle of fx & running http server (Gin)
func RegisterHooks(lifecycle fx.Lifecycle, logger *helper.LogHandler) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Starting application")
				return nil
			},
			OnStop: func(context.Context) error {
				logger.Info("Stopping application")
				return nil
			},
		},
	)
}
