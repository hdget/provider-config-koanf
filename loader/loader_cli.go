package loader

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

type cliConfigLoader struct {
	reader        *koanf.Koanf
	configContent []byte
}

func NewCliConfigLoader(reader *koanf.Koanf, configContent []byte) Loader {
	return &cliConfigLoader{
		reader:        reader,
		configContent: configContent,
	}
}

func (l *cliConfigLoader) Load() error {
	if len(l.configContent) > 0 {
		return l.reader.Load(rawbytes.Provider(l.configContent), toml.Parser())
	}
	return nil
}
