package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "dotger",
		Authors: []*cli.Author{
			{
				Name: "4strodev",
			},
		},
		Usage: "dotger",
		Description: "A dotfiles manager inspired by stow",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Execute main action")
			return nil
		},
	}
	app.Run(os.Args)
}
