package commands

import (
	"strings"
	"tomoribot-geminiai-version/src/commands/instagram"
	"tomoribot-geminiai-version/src/commands/play"
	"tomoribot-geminiai-version/src/commands/shazam"
	"tomoribot-geminiai-version/src/commands/sticker"
	"tomoribot-geminiai-version/src/commands/sticker2"
	"tomoribot-geminiai-version/src/commands/tiktok"
	"tomoribot-geminiai-version/src/commands/tomp3"
	"tomoribot-geminiai-version/src/commands/twitter"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/commands/ytmp3"
	"tomoribot-geminiai-version/src/commands/ytmp4"
)

var Commands = map[string]command_types.Command{
	"play": {Execute: play.Execute, Details: play.Details()},
	"ytmp3": {Execute: ytmp3.Execute, Details: ytmp3.Details()},
	"ytmp4": {Execute: ytmp4.Execute, Details: ytmp4.Details()},
	"instagram": {Execute: instagram.Execute, Details: instagram.Details()},
	"shazam": {Execute: shazam.Execute, Details: shazam.Details()},
	"tomp3": {Execute: tomp3.Execute, Details: tomp3.Details()},
	"sticker": {Execute: sticker.Execute, Details: sticker.Details()},
	"sticker2": {Execute: sticker2.Execute, Details: sticker2.Details()},
	"twitter": {Execute: twitter.Execute, Details: twitter.Details()},
	"tiktok": {Execute: tiktok.Execute, Details: tiktok.Details()},
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
