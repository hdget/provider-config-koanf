package koanf

import (
	"os"

	"github.com/hdget/common/constant"
	"github.com/hdget/common/types"
	"github.com/hdget/provider-config-koanf/loader"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
)

// koanfConfigProvider 命令行配置
type koanfConfigProvider struct {
	reader        *koanf.Koanf
	app           string
	env           string
	configFile    string // 指定的配置文件
	configContent []byte // 指定的配置内容

}

const (
	defaultStructTag = "mapstructure"
)

// New 初始化config provider
func New(app string, options ...Option) (types.ConfigProvider, error) {
	if app == "" {
		return nil, errors.New("app is required")
	}

	provider := &koanfConfigProvider{
		app:    app,
		env:    os.Getenv(constant.EnvKeyRunEnvironment),
		reader: koanf.New("."),
	}

	for _, option := range options {
		option(provider)
	}

	err := provider.Load()
	if err != nil {
		return nil, errors.Wrap(err, "load config")
	}

	return provider, nil
}

// Load 从各个配置源获取配置数据, 并加载到configVar中，同名变量优先级高的覆盖低的
// - configFile: 文件配置(低)
// - input: 命令行参数配置(最高)
// - env: 环境变量配置(高)
func (p *koanfConfigProvider) Load() error {
	// minimal config
	if err := loader.NewMinimalConfigLoader(p.reader, p.app).Load(); err != nil {
		return errors.Wrap(err, "load minimal config")
	}

	if err := loader.NewFileConfigLoader(p.reader, p.app, p.env, p.configFile).Load(); err != nil {
		return errors.Wrap(err, "load config from file")
	}

	if err := loader.NewCliConfigLoader(p.reader, p.configContent).Load(); err != nil {
		return errors.Wrap(err, "load config from cli")
	}

	if err := loader.NewEnvConfigLoader(p.reader).Load(); err != nil {
		return errors.Wrap(err, "load config from env")
	}

	return nil
}

// Unmarshal 解析配置
func (p *koanfConfigProvider) Unmarshal(configVar any, args ...string) error {
	if len(args) > 0 {
		return p.reader.UnmarshalWithConf(args[0], configVar, koanf.UnmarshalConf{Tag: defaultStructTag})
	}
	return p.reader.UnmarshalWithConf("", configVar, koanf.UnmarshalConf{Tag: defaultStructTag})
}

func (p *koanfConfigProvider) Get(key string) any {
	return p.reader.Get(key)
}

func (p *koanfConfigProvider) GetCapability() types.Capability {
	return Capability
}
