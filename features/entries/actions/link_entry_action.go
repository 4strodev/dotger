package actions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/4strodev/dotger/features/entries/domain"
	"github.com/4strodev/dotger/shared/injector"
	"github.com/spf13/afero"
)

func LinkEntry(inject injector.Injector, ctx context.Context, source string, entry domain.EntryConfig) error {
	fs := inject.FileSystem.GetFs()
	exists, err := inject.FileSystem.Exists(entry.Destination.Path)
	if err != nil {
		return err
	}
	if !exists {
		if !entry.Destination.Mkdir {
			return fmt.Errorf("'%s' does not exist", entry.Destination.Path)
		} else {
			err = fs.Mkdir(entry.Destination.Path, os.ModeDir | os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	stats, err := fs.Stat(source)
	if err != nil {
		return err
	}
	if !stats.IsDir() {
		return fmt.Errorf("%s is not a directory", source)
	}

	files, err := afero.ReadDir(fs, source)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == ".dotger.toml" {
			continue
		}
		origin := filepath.Join(source, file.Name())
		target := filepath.Join(entry.Destination.Path, file.Name())
		err := inject.FileSystem.Symlink(origin, target)
		if err != nil {
			return err
		}
	}

	return nil
}
