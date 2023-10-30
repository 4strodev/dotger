package actions

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/4strodev/dotger/features/entries/domain"
	"github.com/4strodev/dotger/shared/injector"
	"github.com/spf13/afero"
)

func LinkEntryAction(inject injector.Injector, ctx context.Context, entry domain.Entry) error {
	fs := inject.FileSystem.GetFs()
	stats, err := fs.Stat(entry.Source)
	if err != nil {
		return err
	}
	if !stats.IsDir() {
		return fmt.Errorf("%s is not a directory", entry.Source)
	}


	files, err := afero.ReadDir(fs, entry.Source)
	if err != nil {
		return err
	}

	for _, file := range files {
		origin := filepath.Join(entry.Source, file.Name())
		target := filepath.Join(entry.Destination, file.Name())
		err := inject.FileSystem.Symlink(origin, target)
		if err != nil {
			return err
		}
	}

	return nil
}
