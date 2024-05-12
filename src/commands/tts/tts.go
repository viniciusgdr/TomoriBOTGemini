package tts

import (
	"math/rand"
	"os"
	"strconv"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	infra_whatsmeow_utils "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/utils"

	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name: "tts",
		Description: "Texto para Áudio",
		Category: command_types.CategoryUtility,
		Permission: command_types.PermissionAll,
		OnlyGroups: true,
		OnlyPrivate: false,
		BotRequiresAdmin: false,
		Alias: []string{"t2a", "text2audio"},
	}
}

func Execute(commandProps *command_types.CommandProps) {
	if len(commandProps.Arg) == 0 {
		return
	}
	
	random := strconv.Itoa(rand.Intn(1000000000))
	speech := htgotts.Speech{Folder: "assets/temp", Language: voices.Portuguese}
	speech.CreateSpeechFile(commandProps.Arg, random)

	buffer, err := os.ReadFile("./assets/temp/" + random + ".mp3")
	if err != nil {
		return
	}
	os.Remove("./assets/temp/" + random + ".mp3")
	device := infra_whatsmeow_utils.GetDevice(commandProps.Message.Info.ID)
	mimetype := infra_whatsmeow_utils.GetMimeTypeAudioByDevice(device)
	sender.SendAudioMessage(
		commandProps.Client.Client,
		commandProps.Message.Info.Chat,
		buffer,
		&sender.MessageOptions{
			MimeType: mimetype,
			QuotedMessage: commandProps.Message,
			Ptt: true,
		},
	)
}