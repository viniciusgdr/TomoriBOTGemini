package ytmp3

import (
	"fmt"
	"strings"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	playServices "tomoribot-geminiai-version/src/services/play"
	web_functions "tomoribot-geminiai-version/src/utils/web"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)
func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name: "ytmp3",
		Description: "Baixar músicas do Youtube em MP3",
		Category: command_types.CategoryDownload,
		Permission: command_types.PermissionAll,
		OnlyGroups: true,
		OnlyPrivate: false,
		BotRequiresAdmin: false,
		Alias: []string{"mp3"},
	}
}


func Execute(commandProps *command_types.CommandProps) {
	if commandProps.Arg == "" {
		commandProps.Reply("É necessário enviar o link do vídeo do YouTube, exemplo: /ytmp3 https://www.youtube.com/watch?v=QH2-TGUlwu4.\n\nCaso queira que o bot envie em documento, adicione --document no final do link.")
		return
	}
	modeDocument := strings.Contains(commandProps.Arg, "--document")
	if modeDocument {
		commandProps.Arg = strings.ReplaceAll(commandProps.Arg, "--document", "")
		commandProps.Arg = strings.Trim(commandProps.Arg, " ")
	}
	id, errVideoId := playServices.GetVideoID(commandProps.Arg)
	if errVideoId != nil {
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
	contentAudio, errAudio := streamings.GetHighAudio()
	if errAudio != nil {
		commandProps.Reply("Algo de errado aconteceu ao procurar o audio do conteúdo.")
		return
	}
	buffer, sizeFile, errDownload := web_functions.GetBufferFromUrlThreads(contentAudio.Url)
	fmt.Println(errDownload)
	if errDownload != nil {
		commandProps.Reply("Ocorreu um erro ao baixar a música, tentando novamente em alguns segundos...")
		return
	}
	if modeDocument || sizeFile > 15728640 {
		sender.SendDocumentMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			info.Title+`.mp3`,
			`• Titulo: + `+info.Title+`
• Canal: `+info.Author+``,
			buffer,
			&sender.MessageOptions{
				MimeType:      "audio/mpeg",
				QuotedMessage: commandProps.Message,
			},
		)
		commandProps.React("✅")
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
					Title:                 proto.String(info.Title),
					MediaType:             waProto.ContextInfo_ExternalAdReplyInfo_VIDEO.Enum(),
					ThumbnailURL:          proto.String(`https://i.ytimg.com/vi/` + id + `/0.jpg`),
					SourceURL:             proto.String("https://www.youtube.com/watch?v=" + id),
					MediaURL:              proto.String("https://www.youtube.com/watch?v=" + id),
					ShowAdAttribution:     proto.Bool(true),
					ContainsAutoReply:     proto.Bool(true),
					RenderLargerThumbnail: proto.Bool(true),
				},
				MentionedJID: []string{
					commandProps.Message.Info.Sender.ToNonAD().String(),
				},
			},
		},
	)
	commandProps.React("✅")
}