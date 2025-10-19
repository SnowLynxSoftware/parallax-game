package cmd

import (
	"errors"
	"os"
)

type Command interface {
	Execute() error
}

type Handler struct {
	commands map[string]Command
}

func NewHandler() *Handler {

	cmds := map[string]Command{
		"server":  &ServerCommand{},
		"migrate": &MigrateCommand{},
	}

	return &Handler{
		commands: cmds,
	}
}

func (h *Handler) ExecuteCommand() error {
	args := os.Args

	var command string
	if len(args) > 1 {
		command = args[1]
	}

	if command == "" {
		command = "server"
	}

	cmd, commandExists := h.commands[command]

	if commandExists {
		err := cmd.Execute()
		return err
	}
	return errors.New("command not found")
}
