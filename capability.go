package koanf

import (
	"github.com/hdget/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "config-koanf"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryConfig,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}
