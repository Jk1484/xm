package gateways

import (
	"xm/gateways/nats"

	"go.uber.org/fx"
)

var Module = fx.Options(
	nats.Module,
)
