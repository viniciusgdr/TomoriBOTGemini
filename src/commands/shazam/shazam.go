package shazam

import (
	"os"
	"tomoribot-geminiai-version/src/commands/play"
	"tomoribot-geminiai-version/src/commands/tomp3"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	infra_whatsmeow_utils "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/utils"
	shazamService "tomoribot-geminiai-version/src/services/shazam"
	web_functions "tomoribot-geminiai-version/src/utils/web"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "shazam",
		Description:      "Identificar músicas via áudio que você não conhece",
		Category:         command_types.CategorySearch,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       true,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
		Alias:            []string{"musica"},
	}
}

func Execute(commandProps *command_types.CommandProps) {
	audio := commandProps.Message.Message.GetAudioMessage()
	quotedMsgAudio := commandProps.QuotedMsg.GetAudioMessage()
	quotedMsgVideo := commandProps.QuotedMsg.GetVideoMessage()
	video := commandProps.Message.Message.GetVideoMessage()

	if audio == nil && quotedMsgAudio == nil && quotedMsgVideo == nil && video == nil {
		commandProps.Reply("Você precisa enviar um áudio ou video para eu identificar a música.")
		return
	}

	mediaByte := []byte{}
	switch {
	case audio != nil:
		mediaByte, _ = commandProps.Client.Client.Download(audio)
	case quotedMsgAudio != nil:
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgAudio)
	case quotedMsgVideo != nil:
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgVideo)
		go commandProps.Reply("Aguarde um momento, estou convertendo o áudio...")
		mediaByteConv, err2 := tomp3.Mp4ToMp3(mediaByte, "mp4")
		if err2 != nil {
			commandProps.Reply("Não foi possível converter o áudio.")
			return
		}
		mediaByte = mediaByteConv
	case video != nil:
		mediaByte, _ = commandProps.Client.Client.Download(video)
		go commandProps.Reply("Aguarde um momento, estou convertendo o áudio...")
		mediaByteConv, err2 := tomp3.Mp4ToMp3(mediaByte, "mp4")
		if err2 != nil {
			commandProps.Reply("Não foi possível converter o áudio.")
			return
		}
		mediaByte = mediaByteConv
	}

	if len(mediaByte) == 0 {
		commandProps.Reply("Não foi possível baixar o áudio.")
		return
	}
	tempPath := infra_whatsmeow_utils.GenerateTempFileName(".ogg")
	err := os.WriteFile(tempPath, mediaByte, 0644)

	if err != nil {
		commandProps.Reply("Não foi possível salvar o áudio.")
		return
	}

	commandProps.Reply("Aguarde um momento, estou identificando a música...")
	shazamResult, err := shazamService.ShazamService(tempPath)
	os.Remove(tempPath)
	if err != nil {
		commandProps.Reply("Não foi possível identificar a música.")
		return
	}
	var lyrics string = ""
	if shazamResult.Track.Sections == nil {
		commandProps.Reply("Não foi possível identificar a música. Aconteceu um erro ou não conseguimos encontrar a música.")
		return
	}

	for _, section := range shazamResult.Track.Sections.([]interface{}) {
		sectionMap := section.(map[string]interface{})
		sectionType, _ := sectionMap["type"].(string)
		sectionText, textExists := sectionMap["text"].([]interface{})
		if textExists && sectionType == "LYRICS" {
			for _, line := range sectionText {
				lyrics += line.(string) + "\n"
			}
			break
		}
	}
	var platforms string = ""
	for _, providers := range shazamResult.Track.Hub.Providers {
		platforms += providers.Type + ": " + providers.Actions[0].URI + "\n"
	}
	image := shazamResult.Track.Images.Background
	buffer, err := web_functions.GetBufferFromUrl(image)
	commandProps.Arg = shazamResult.Track.Title + " " + shazamResult.Track.Subtitle
	if err != nil {
		commandProps.Reply(`🔍 Provavelmente esse áudio contém a música *` + shazamResult.Track.Title + `* de *` + shazamResult.Track.Subtitle + `*. 🤓`)
		if shazamResult.Track.Title != "" {
			play.Execute(commandProps)
		}
		return
	}
	sender.SendImageMessage(
		commandProps.Client.Client,
		commandProps.Message.Info.Chat,
		`🔍 Provavelmente esse áudio contém a música *`+shazamResult.Track.Title+`* de *`+shazamResult.Track.Subtitle+`*. 🤓`,
		buffer,
		&sender.MessageOptions{})
	if shazamResult.Track.Title != "" {
		play.Execute(commandProps)
	}
}
