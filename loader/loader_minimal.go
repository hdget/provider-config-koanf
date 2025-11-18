package loader

import (
	"fmt"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

const (
	// 最小化的配置,保证日志工作正常
	tplMinimalConfigContent = `
[sdk]
  [sdk.logger]
	   level = "debug"
	   filename = "%s.log"
	   [sdk.logger.rotate]
		   max_age = 7`
)

type minimalConfigLoader struct {
	app    string
	reader *koanf.Koanf
}

func NewMinimalConfigLoader(reader *koanf.Koanf, app string) Loader {
	return &minimalConfigLoader{
		app:    app,
		reader: reader,
	}
}

func (l *minimalConfigLoader) Load() error {
	minimalConfig := []byte(fmt.Sprintf(tplMinimalConfigContent, l.app))
	return l.reader.Load(rawbytes.Provider(minimalConfig), toml.Parser())
}
