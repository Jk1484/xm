package repositories

import (
	"xm/pkg/repositories/company"
	"xm/pkg/repositories/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	user.Module,
	company.Module,
)
