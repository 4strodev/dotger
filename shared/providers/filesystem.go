package providers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	promise "github.com/4strodev/promise/pkg"
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
func (fs *FileSystem) CopyDir(ctx context.Context, origin string, destination string) *promise.Promise[struct{}] {
	return promise.New(func(resolve func(struct{}), reject func(error)) {
		// Checking if provided paths exists
		originStat, err := fs.fs.Stat(origin)
		if err != nil {
			reject(err)
			return
		}

		if fs.IsSymlink(origin) {
			target, err := fs.ReadLink(origin)
			if err != nil {
				reject(err)
				return
			}

			err = fs.Symlink(target, destination)
			if err != nil {
				reject(err)
				return
			}

			resolve(struct{}{})
			return
		}

		// Getting destination status
		_, err = fs.fs.Stat(destination)
		if err == nil {
			reject(fmt.Errorf("%s already exists", destination))
			return
		}

		// If cannot get status for some other reason
		// rather than directory does not exist then
		// reutrn the error
		if !os.IsNotExist(err) {
			reject(err)
			return
		}

		// If destination directory does not exits
		// create a new one
		err = fs.fs.Mkdir(destination, os.ModePerm)
		if err != nil {
			reject(err)
			return
		}

		if !originStat.IsDir() {
			reject(fmt.Errorf("%s is not a directory", origin))
			return
		}

		// Reading files of origin
		files, err := afero.ReadDir(fs.fs, origin)
		if err != nil {
			reject(err)
			return
		}

		// Copying files
		promises := make([]*promise.Promise[struct{}], 0)
		for _, file := range files {
			originName := path.Join(origin, file.Name())
			destinationName := path.Join(destination, file.Name())
			if file.IsDir() {
				prom := fs.CopyDir(ctx, originName, destinationName)
				promises = append(promises, prom)
			} else {
				prom := fs.CopyFile(originName, destinationName)
				promises = append(promises, prom)
			}
		}

		_, err = promise.MergeAll(ctx, promises...).Await(ctx)
		if err != nil {
			reject(err)
			return
		}
		resolve(struct{}{})
	})
}

// CopyFile copies origin file to destination path
func (fs *FileSystem) CopyFile(origin string, destination string) *promise.Promise[struct{}] {
	return promise.New(func(resolve func(struct{}), reject func(error)) {
		// Checking if origin file exists
		originStat, err := fs.fs.Stat(origin)
		if err != nil {
			reject(err)
			return
		}

		// If origin is a symlink then create a new symlink
		if fs.IsSymlink(origin) {
			target, err := fs.ReadLink(origin)
			if err != nil {
				reject(err)
				return
			}
			err = fs.Symlink(target, destination)
			if err != nil {
				reject(err)
				return
			}
			resolve(struct{}{})
			return
		}

		// If origin is a regular file then copy content
		// Opening origin file
		originFile, err := fs.fs.Open(origin)
		if err != nil {
			reject(err)
			return
		}
		defer originFile.Close()

		// Creating detination file
		destinationFile, err := fs.fs.OpenFile(destination, os.O_CREATE|os.O_RDWR|os.O_TRUNC, originStat.Mode())
		defer destinationFile.Close()
		if err != nil {
			reject(err)
			return
		}

		// Reading origin file content
		content, err := afero.ReadFile(fs.fs, origin)
		if err != nil {
			reject(err)
			return
		}
		// Writing content to destination file
		_, err = destinationFile.Write(content)
		if err != nil {
			reject(err)
			return
		}

		resolve(struct{}{})
	})
}
