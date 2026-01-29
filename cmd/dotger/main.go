package main

import (
	"log"
	"os"

	"github.com/4strodev/dotger/features/entries/infrastructure"
	"github.com/4strodev/dotger/shared/injector"
	"github.com/urfave/cli/v2"
)

func main() {
	inject := injector.NewInjector()
	app := &cli.App{
		Name: "dotger",
		Authors: []*cli.Author{
			{
				Name: "4strodev",
			},
		},
		Usage: "Centralize and manage your dotfiles",
		Description: "A dotfiles manager inspired by stow",
		Commands: []*cli.Command{
			infrastructure.GetLinkCommand(inject),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
