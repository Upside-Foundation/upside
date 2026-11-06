package main

import (
	"os"

	"github.com/0xUp5ide/upside-foundation/app"
	"github.com/0xUp5ide/upside-foundation/cmd/upsided/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, cmd.EnvPrefix, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
