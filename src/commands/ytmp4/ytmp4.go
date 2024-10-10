package ytmp4

import (
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	playServices "tomoribot-geminiai-version/src/services/play"
	web_functions "tomoribot-geminiai-version/src/utils/web"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "ytmp4",
		Description:      "Baixar músicas do Youtube em MP4",
		Category:         command_types.CategoryDownload,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       true,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
		Alias:            []string{"mp4"},
	}
}

func Execute(commandProps *command_types.CommandProps) {
	if commandProps.Arg == "" {
		commandProps.Reply("É necessário enviar o link do vídeo do YouTube, exemplo: /ytmp4 https://www.youtube.com/watch?v=QH2-TGUlwu4")
		return
	}
	id, err := playServices.GetVideoID(commandProps.Arg)
	if err != nil {
		result, err := playServices.Search(commandProps.Arg)
		if len(result) == 0 || err != nil {
			commandProps.Reply("Não encontrei nenhuma música com esse nome")
			return
		}
		id = result[0].VideoID
	}
	go commandProps.React("🔎")
	info, streamings, errVideoInfo := playServices.GetVideoInfo(id)
	if errVideoInfo != nil {
		commandProps.Reply("Ocorreu um erro ao procurar a música.")
		return
	}
	contentVideo, errVideo := streamings.GetHighVideo()
	if errVideo != nil {
		commandProps.Reply("Algo de errado aconteceu ao procurar o vídeo. Tente novamente mais tarde.")
		return
	}
	buffer, sizeFile, errDownload := web_functions.GetBufferFromUrlThreads(contentVideo.Url)
	if errDownload != nil {
		commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
	}
	if sizeFile > 1.5 *1024*1024*1024 {
		commandProps.Reply("O vídeo é muito grande para ser enviado.")
		return
	}

	if sizeFile > 15*1024*1024 {
		sender.SendDocumentMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			info.Title+`.mp4`,
			`• Titulo: `+info.Title+`
• Canal: `+info.Author+`
• Qualidade: `+contentVideo.Quality+``,
			buffer,
			&sender.MessageOptions{
				MimeType:      "video/mp4",
				QuotedMessage: commandProps.Message,
			},
		)
	} else {
		sender.SendVideoMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			`• Titulo: `+info.Title+`
• Canal: `+info.Author+``,
			buffer,
			&sender.MessageOptions{
				MimeType:      "video/mp4",
				QuotedMessage: commandProps.Message,
			},
		)
	}
	commandProps.React("✅")
}
