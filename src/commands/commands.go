package commands

import (
	"strings"
	"tomoribot-geminiai-version/src/commands/play"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/commands/ytmp3"
	"tomoribot-geminiai-version/src/commands/ytmp4"
)

var Commands = map[string]command_types.Command{
	"play": {Execute: play.Execute, Details: play.Details()},
	"ytmp3": {Execute: ytmp3.Execute, Details: ytmp3.Details()},
	"ytmp4": {Execute: ytmp4.Execute, Details: ytmp4.Details()},
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
	commandNameLower := strings.ToLower(commandName)
	command := command_types.Command{}
	for _, cmd := range AllCommands {
		if cmd.Details.Name == commandNameLower {
			command = cmd
		}
		for _, alias := range cmd.Details.Alias {
			if alias == commandNameLower {
				command = cmd
			}
		}
	}
	return command, command.Details.Name != ""
}
