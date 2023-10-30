package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/4strodev/dotger/features/entries/actions"
	"github.com/4strodev/dotger/features/entries/domain"
	"github.com/4strodev/dotger/shared/injector"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/urfave/cli/v2"
)

func GetLinkCommand(inject injector.Injector) *cli.Command {
	return &cli.Command {
		Name: "link",
		Usage: "link a dotger entry to their destination",
		Action: func(ctx *cli.Context) error {
			const configFile = ".dotger.toml"
			args := ctx.Args()
			if args.Len() != 1 {
				return errors.New("Expected 1 argument")
			}

			entryName := args.First()
			exists, err := inject.FileSystem.Exists(entryName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("entry '%s' doesn't exists", entryName)
			}

			dotgerFile := filepath.Join(entryName, configFile)
			exists, err = inject.FileSystem.Exists(dotgerFile)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("entry '%s' doesn't have '%s' file", entryName, configFile)
			}

			content, err := inject.FileSystem.ReadFile(dotgerFile)
			if err != nil {
				return err
			}

			var entryConfig domain.EntryConfig
			k := koanf.New(".")
			k.Load(rawbytes.Provider(content), toml.Parser())
			err = k.Unmarshal("", &entryConfig)
			if err != nil {
				return err
			}

			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			err = actions.LinkEntry(inject, ctx.Context, filepath.Join(wd, entryName), entryConfig)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
