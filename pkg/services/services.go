package services

import (
	"xm/pkg/services/company"
	"xm/pkg/services/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	user.Module,
	company.Module,
)
