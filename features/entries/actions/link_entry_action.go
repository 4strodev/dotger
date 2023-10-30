package actions

import (
	"context"
	"fmt"

	"github.com/4strodev/dotger/features/entries/domain"
	"github.com/4strodev/dotger/shared/injector"
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

	err = inject.FileSystem.Symlink(entry.Source, entry.Destination)
	if err != nil {
		return err
	}

	return nil
}
