package koanf

type Option func(provider *koanfConfigProvider)

func WithConfigFile(configFile string) Option {
	return func(p *koanfConfigProvider) {
		p.configFile = configFile
	}
}

func WithConfigContent(configContent []byte) Option {
	return func(p *koanfConfigProvider) {
		p.configContent = configContent
	}
}
