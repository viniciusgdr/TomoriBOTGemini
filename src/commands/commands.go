package commands

import (
	command_types "tomoribot-geminiai-version/src/commands/types"
)

var Commands = map[string]command_types.Command{
}

func LoadCommands() []command_types.Command {
	var commands []command_types.Command
	for _, commandContent := range Commands {
		commands = append(commands, command_types.Command{
			Execute: commandContent.Execute,
			Details: commandContent.Details,
		})
	}
	return commands
}

var AllCommands = LoadCommands()

func GetCommand(commandName string) (command_types.Command, bool) {
	command := command_types.Command{}
	for _, cmd := range AllCommands {
		if cmd.Details.Name == commandName {
			command = cmd
		}
		for _, alias := range cmd.Details.Alias {
			if alias == commandName {
				command = cmd
			}
		}
	}
	return command, command.Details.Name != ""
}
