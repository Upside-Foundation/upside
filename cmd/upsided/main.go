package main

import (
	"os"

	"github.com/Upside-Foundation/upside/app"
	"github.com/Upside-Foundation/upside/cmd/upsided/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, cmd.EnvPrefix, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
