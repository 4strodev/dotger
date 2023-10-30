package injector

import (
	"github.com/4strodev/dotger/shared/providers"
	"github.com/spf13/afero"
)

type Injector struct {
	FileSystem *providers.FileSystem
}

func NewInjector() Injector {
	fs := providers.NewFileSystem(afero.NewOsFs())
	return Injector{
		FileSystem: fs,
	}
}

