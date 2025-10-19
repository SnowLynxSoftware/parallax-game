package cmd

import (
	"github.com/morpheuszero/go-heimdall"
	"github.com/snowlynxsoftware/parallax-game/config"
)

type MigrateCommand struct {
}

func (s *MigrateCommand) Execute() error {
	appConfig := config.NewAppConfig()
	h := heimdall.NewHeimdall(appConfig.GetDBConnectionString(), "migration_history", "./migrations", true)
	return h.RunMigrations()
}
