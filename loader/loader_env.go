package loader

import (
	"strings"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

type envLoader struct {
	reader *koanf.Koanf
}

const (
	defaultEnvPrefix = "HD_"
)

func NewEnvConfigLoader(reader *koanf.Koanf) Loader {
	return &envLoader{
		reader: reader,
	}
}

// Load 从环境变量中读取配置信息
func (l *envLoader) Load() error {
	// Load only environment variables with prefix "HD_" and merge into config.
	// Transform var names by:
	// 1. Converting to lowercase
	// 2. Removing "HD_" prefix
	// 3. Replacing "_" with "." to representing nesting using the . delimiter.
	// Example: HD_PARENT1_CHILD1_NAME becomes "parent1.child1.name"
	return l.reader.Load(env.Provider(".", env.Opt{
		Prefix: defaultEnvPrefix,
		TransformFunc: func(k, v string) (string, any) {
			// Transform the key.
			k = strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, defaultEnvPrefix)), "_", ".")

			// Transform the value into slices, if they contain spaces.
			// e,g: HD_TAGS="foo bar baz" -> tags: ["foo", "bar", "baz"]
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil)
}
