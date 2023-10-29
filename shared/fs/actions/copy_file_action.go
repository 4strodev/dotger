package actions

import "github.com/spf13/afero"

func CopyFileAction(fileSystem afero.Fs, source string, destination string) error {

	stats, err := fileSystem.Stat(source)
	if err != nil {
		return err
	}
	return nil
}
