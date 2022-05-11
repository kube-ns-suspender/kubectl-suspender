package main

import (
	"github.com/govirtuo/kubectl-suspender/app"
	"github.com/govirtuo/kubectl-suspender/cmd"
)

var (
	// those variables are filled at link time
	BinaryName, Version, BuildDate string
)

func main() {
	a, err := app.New()
	if err != nil {
		a.Logger.Fatal().Err(err)
	}
	a.Logger.Info().Msgf("%s version '%s' (built: %s)", BinaryName, Version, BuildDate)
	cmd.Execute(a)
}
