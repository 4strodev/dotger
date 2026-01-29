package infrastructure

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/4strodev/dotger/features/entries/actions"
	"github.com/4strodev/dotger/features/entries/domain"
	"github.com/4strodev/dotger/shared/injector"
	"github.com/urfave/cli/v2"
)

func GetLinkCommand(inject injector.Injector) *cli.Command {
	return &cli.Command{
		Name:      "link",
		ArgsUsage: "<entry path>",
		Usage:     "link a dotger entry to their destination",
		Aliases:   []string{"l"},
		Action: func(ctx *cli.Context) error {
			const configFile = ".dotger.toml"
			// Checcking if command line arguments are correct
			args := ctx.Args()
			if args.Len() != 1 {
				return errors.New("Expected <entry path> argument")
			}

			// Getting entry name from command line arguments
			entryPath := args.First()
			exists, err := inject.FileSystem.Exists(entryPath)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("No dotger entry found at '%s'", entryPath)
			}

			// Getting full path of dotfer file
			dotgerFile := filepath.Join(entryPath, configFile)
			exists, err = inject.FileSystem.Exists(dotgerFile)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("No '%s' file found at '%s'", configFile, entryPath)
			}

			// Reading config file
			content, err := inject.FileSystem.ReadFile(dotgerFile)
			if err != nil {
				return err
			}

			// Unmarshalling config file
			configFileParser := domain.EntryConfigFileParser{}
			entryConfig, err := configFileParser.Parse(string(content))
			if err != nil {
				return err
			}

			// Getting current working directory
			absoluteEntryPath, err := filepath.Abs(entryPath)
			if err != nil {
				return err
			}
			err = actions.LinkEntry(inject, ctx.Context, absoluteEntryPath, entryConfig)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
