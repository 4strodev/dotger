package providers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSymlink(t *testing.T) {
	const message string = "Hello world!"
	const original string = "original.txt"
	const symlink string = "new.txt"
	fileSystem := NewFileSystem(afero.NewOsFs())

	// Creating original file
	file, err := fileSystem.GetFs().Create(original)
	assert.NoError(t, err)
	_, err = file.WriteString(message)
	assert.NoError(t, err)
	defer file.Close()
	defer fileSystem.GetFs().Remove(original)

	// Creaing symlink
	err = fileSystem.Symlink(original, symlink)
	assert.NoError(t, err)

	// Opening symlink 
	newFile, err := fileSystem.GetFs().Open(symlink)
	assert.NoError(t, err)
	defer newFile.Close()
	defer fileSystem.GetFs().Remove(symlink)

	// Checking if the content is the same
	content, err := afero.ReadAll(newFile)
	assert.NoError(t, err)
	assert.Equal(t, message, string(content))
}

func TestCopyFile(t *testing.T) {
	fileSystem := NewFileSystem(afero.NewMemMapFs())
	const message string = "Hello world!"
	const original string = "original.txt"
	const destination string = "new.txt"

	file, err := fileSystem.GetFs().Create(original)
	assert.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString(message)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = fileSystem.CopyFile(original, destination).Await(ctx)
	assert.NoError(t, err)
	
	content, err := afero.ReadFile(fileSystem.GetFs(), destination)
	assert.NoError(t, err)

	assert.Equal(t, string(content), message)
}

func TestCopySymlink(t *testing.T) {
	fileSystem := NewFileSystem(afero.NewOsFs())
	const message string = "Hello world!"
	const original string = "original.txt"
	const symlink string = "link_original.txt"
	const destination string = "new.txt"

	// Creating original file
	file, err := fileSystem.GetFs().Create(original)
	assert.NoError(t, err)
	defer file.Close()
	defer fileSystem.GetFs().Remove(original)

	// Writing content to file
	_, err = file.WriteString(message)
	assert.NoError(t, err)

	// Creating symlink
	err = fileSystem.Symlink(original, symlink)
	assert.NoError(t, err)
	defer fileSystem.GetFs().Remove(symlink)

	// Copying file
	ctx := context.Background()
	_, err = fileSystem.CopyFile(symlink, destination).Await(ctx)
	assert.NoError(t, err)
	defer fileSystem.GetFs().Remove(destination)
	
	content, err := afero.ReadFile(fileSystem.GetFs(), destination)
	assert.NoError(t, err)

	assert.Equal(t, string(content), message)
	assert.True(t, fileSystem.IsSymlink(destination))
}
