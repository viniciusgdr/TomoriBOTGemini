package play

import (
	"fmt"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/commands/ytmp3"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	playServices "tomoribot-geminiai-version/src/services/play"
	web_functions "tomoribot-geminiai-version/src/utils/web"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)
func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name: "play",
		Description: "Baixar músicas do Youtube em MP3 por Busca",
		Category: command_types.CategoryDownload,
		Permission: command_types.PermissionAll,
		OnlyGroups: false,
		OnlyPrivate: false,
		BotRequiresAdmin: false,
		Alias: []string{"p"},
	}
}

func Execute(commandProps *command_types.CommandProps) {
	if commandProps.Arg == "" {
		commandProps.Reply("Insira o nome da música logo após o comando, exemplo: /play Proteção de Tela - Tarcísio do Acordeon")
		return
	}
	id, err := playServices.GetVideoID(commandProps.Arg)
	if len(id) > 0 && err == nil {
		ytmp3.Execute(commandProps)
		return
	}
	go commandProps.React("🔎")
	result, err := playServices.Search(commandProps.Arg)
	if len(result) == 0 || err != nil {
		commandProps.Reply("Não encontrei nenhuma música com esse nome")
		return
	}
	success := false
	lengthRetrys := 0
	for !success {
		if lengthRetrys >= len(result) {
			commandProps.Reply("Ocorreu um erro ao baixar a música, tente novamente mais tarde")
			return
		}
		fmt.Println("Getting Music Video Info from", result[lengthRetrys].VideoID)
		videoInfo, streamings, err2 := playServices.GetVideoInfo(result[lengthRetrys].VideoID)
		if err2 != nil {
			lengthRetrys++
			continue
		}
		streaming, err3 := streamings.GetHighAudio()
		if err3 != nil {
			lengthRetrys++
			continue
		}
		success = true
		buffer, sizeFile, err4 := web_functions.GetBufferFromUrlThreads(streaming.Url)
		if err4 != nil {
			lengthRetrys++
			success = false
			commandProps.Reply("Ocorreu um erro ao baixar a música, tentando novamente em alguns segundos...")
			continue
		}
		if sizeFile > 15728640 {
			sender.SendDocumentMessage(
				commandProps.Client.Client,
				commandProps.Message.Info.Chat,
				videoInfo.Title + ".mp3",
				"",
				buffer,
				&sender.MessageOptions{
					MimeType: "audio/mpeg",
					ContextInfo: &waProto.ContextInfo{
						MentionedJid: []string{
							commandProps.Message.Info.Sender.ToNonAD().String(),
						},
					},
				},
			)
			go commandProps.React("✅")
			return
		}
		sender.SendAudioMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			buffer,
			&sender.MessageOptions{
				MimeType: "audio/mpeg",
				ContextInfo: &waProto.ContextInfo{
					ExternalAdReply: &waProto.ContextInfo_ExternalAdReplyInfo{
						Title:                 proto.String(videoInfo.Title),
						MediaType:             waProto.ContextInfo_ExternalAdReplyInfo_VIDEO.Enum(),
						ThumbnailUrl:          proto.String(`https://i.ytimg.com/vi/` + result[lengthRetrys].VideoID + `/0.jpg`),
						SourceUrl:             proto.String("https://www.youtube.com/watch?v=" + result[lengthRetrys].VideoID),
						MediaUrl:              proto.String("https://www.youtube.com/watch?v=" + result[lengthRetrys].VideoID),
						ShowAdAttribution:     proto.Bool(true),
						ContainsAutoReply:     proto.Bool(true),
						RenderLargerThumbnail: proto.Bool(true),
					},
					MentionedJid: []string{
						commandProps.Message.Info.Sender.ToNonAD().String(),
					},
				},
			},
		)
		commandProps.React("✅")
	}
}