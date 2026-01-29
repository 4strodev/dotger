package domain

type EntryConfig struct {
	Destination struct {
		Path  string `koanf:"path"`
		Mkdir bool   `koanf:"mkdir"`
	} `koanf:"destination"`
}

