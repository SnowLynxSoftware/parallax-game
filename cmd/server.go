package cmd

import (
	"github.com/snowlynxsoftware/parallax-game/config"
	"github.com/snowlynxsoftware/parallax-game/server"
)

type ServerCommand struct {
}

func (s *ServerCommand) Execute() error {
	appConfig := config.NewAppConfig()

	server := server.NewAppServer(appConfig)
	server.Start()
	return nil
}
