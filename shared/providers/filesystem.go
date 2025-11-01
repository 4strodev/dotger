package providers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/spf13/afero"
)

type FileSystem struct {
	fs afero.Fs
}

func NewFileSystem(fs afero.Fs) *FileSystem {
	return &FileSystem{
		fs,
	}
}

// Checks if a file, directory or symlink exists
func (fs *FileSystem) Exists(path string) (bool, error) {
	return afero.Exists(fs.fs, path)
}

// Read a file and return their content. If the file is a symlink it
// follows the link to the original file
func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	var content []byte
	file, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	content, err = afero.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Checks if provided path is a symlink
// if the underlying file system does not admit
// symlinks it will return false
func (fs *FileSystem) IsSymlink(target string) bool {
	lstater, ok := fs.fs.(afero.Lstater)
	if !ok {
		return false
	}

	stats, calledLstat, err := lstater.LstatIfPossible(target)
	if err != nil {
		return false
	}
	if !calledLstat {
		return false
	}

	return stats.Mode()&os.ModeSymlink == os.ModeSymlink
}

// Given a path pointing to a symlink returns the target of that symlink
func (fs *FileSystem) ReadLink(target string) (string, error) {
	linkReader, ok := fs.fs.(afero.LinkReader)
	if !ok {
		return "", errors.New("Cannot read symlinks")
	}
	return linkReader.ReadlinkIfPossible(target)
}

// Creates a symlink
func (fs *FileSystem) Symlink(oldname, symlink string) error {
	linker, ok := fs.fs.(afero.Linker)
	if !ok {
		return errors.New("Cannot create symbolic links")
	}

	err := linker.SymlinkIfPossible(oldname, symlink)
	if err != nil {
		return err
	}

	return nil
}

// Return the underlying file system
func (fs *FileSystem) GetFs() afero.Fs {
	return fs.fs
}

// CopyDir copies content of directory origin to destination
// creating destination directory. If the parent directory of
// destination does not exits it will return an error.
// If destination directory already exits, it return an error.
func (fs *FileSystem) CopyDir(ctx context.Context, origin string, destination string) error {
	// Checking if provided paths exists
	originStat, err := fs.fs.Stat(origin)
	if err != nil {
		return err
	}

	if fs.IsSymlink(origin) {
		target, err := fs.ReadLink(origin)
		if err != nil {
			return err
		}

		err = fs.Symlink(target, destination)
		if err != nil {
			return err
		}

		return nil
	}

	// Getting destination status
	_, err = fs.fs.Stat(destination)
	if err == nil {
		return fmt.Errorf("%s already exists", destination)
	}

	// If cannot get status for some other reason
	// rather than directory does not exist then
	// reutrn the error
	if !os.IsNotExist(err) {
		return err
	}

	// If destination directory does not exits
	// create a new one
	err = fs.fs.Mkdir(destination, os.ModePerm)
	if err != nil {
		return err
	}

	if !originStat.IsDir() {
		return fmt.Errorf("%s is not a directory", origin)
	}

	// Reading files of origin
	files, err := afero.ReadDir(fs.fs, origin)
	if err != nil {
		return err
	}

	// Copying files
	wg := sync.WaitGroup{}
	errChannel := make(chan error)

	// Dispatching IO operations
	for _, file := range files {
		wg.Add(1)
		go func(file os.FileInfo, errChannel chan error) {
			defer wg.Done()
			originName := path.Join(origin, file.Name())
			destinationName := path.Join(destination, file.Name())
			var err error
			if file.IsDir() {
				err = fs.CopyDir(ctx, originName, destinationName)
			} else {
				err = fs.CopyFile(originName, destinationName)
			}

			errChannel <- err

		}(file, errChannel)
	}

	go func() {
		wg.Wait()
		close(errChannel)
	}()

	errorList := []error{}
	for receivedError := range errChannel {
		errorList = append(errorList, receivedError)
	}

	return errors.Join(errorList...)
}

// CopyFile copies origin file to destination path
func (fs *FileSystem) CopyFile(origin string, destination string) error {
	// Checking if origin file exists
	originStat, err := fs.fs.Stat(origin)
	if err != nil {
		return err
	}

	// If origin is a symlink then create a new symlink
	if fs.IsSymlink(origin) {
		target, err := fs.ReadLink(origin)
		if err != nil {
			return err
		}
		err = fs.Symlink(target, destination)
		if err != nil {
			return err
		}
		return nil
	}

	// If origin is a regular file then copy content
	// Opening origin file
	originFile, err := fs.fs.Open(origin)
	if err != nil {
		return err
	}
	defer originFile.Close()

	// Creating detination file
	destinationFile, err := fs.fs.OpenFile(destination, os.O_CREATE|os.O_RDWR|os.O_TRUNC, originStat.Mode())
	defer destinationFile.Close()
	if err != nil {
		return err
	}

	// Reading origin file content
	content, err := afero.ReadFile(fs.fs, origin)
	if err != nil {
		return err
	}
	// Writing content to destination file
	_, err = destinationFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}
