package main

import (
	"context"
	"os"

	"github.com/kaankoken/versioning-tool/cmd"
)

func main() {
	args := os.Args[1:]

	app, err := cmd.MainApp(args)

	if err == nil {
		app.Run()

		defer app.Stop(context.Background())
	} else {
		panic(err)
	}
}
