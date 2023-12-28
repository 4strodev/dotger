package domain

import (
	"bytes"
	"os"
	"text/template"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
)

type EntryConfigFileParser struct {
}

func (c EntryConfigFileParser) Parse(content string) (EntryConfig, error) {
	// Unmarshalling config file
	var entryConfig EntryConfig
	k := koanf.New(".")
	k.Load(rawbytes.Provider([]byte(content)), toml.Parser())
	err := k.Unmarshal("", &entryConfig)
	if err != nil {
		return entryConfig, err
	}

	// Content is loaded now parse strings
	templateResult := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("destination.path").Funcs(template.FuncMap{
		"getenv": os.Getenv,
	}).Parse(entryConfig.Destination.Path)
	if err != nil {
		return entryConfig, err
	}
	err = tmpl.
		Execute(templateResult, nil)
	if err != nil {
		return entryConfig, err
	}
	entryConfig.Destination.Path = templateResult.String()

	return entryConfig, nil
}
