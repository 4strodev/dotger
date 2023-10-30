package infrastructure

type EntryConfig struct {
	Destination string `koanf:"destination"`
	Mkdir       bool   `koanf:"mkdir"`
}
