package infrastructure

import (
	"errors"
	"fmt"

	"github.com/4strodev/dotger/shared/injector"
	"github.com/urfave/cli/v2"
)

func GetLinkCommand(inject injector.Injector) *cli.Command {
	return &cli.Command {
		Name: "link",
		Usage: "link a dotger entry to their destination",
		Action: func(ctx *cli.Context) error {
			args := ctx.Args()
			if args.Len() != 1 {
				return errors.New("Expected 1 argument")
			}

			entryName := args.First()
			fmt.Printf("TODO link entry %s\n", entryName)
			return nil
		},
	}
}
