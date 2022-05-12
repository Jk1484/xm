package main

import (
	"xm/configs"
	"xm/gateways"
	"xm/pkg/db"
	"xm/pkg/handlers"
	"xm/pkg/handlers/server"
	"xm/pkg/logger"
	"xm/pkg/repositories"
	"xm/pkg/services"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Options(
			configs.Module,
			logger.Module,
			db.Module,
			repositories.Module,
			services.Module,
			handlers.Module,
			server.Module,
			gateways.Module,
		),
	).Run()
}
