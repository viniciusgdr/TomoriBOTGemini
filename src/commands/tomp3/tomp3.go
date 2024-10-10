package tomp3

import (
	"os"
	"os/exec"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	infra_whatsmeow_utils "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/utils"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "tomp3",
		Description:      "Converter Video para Áudio",
		Category:         command_types.CategoryUtility,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       false,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
		Alias:            []string{"mp3"},
	}
}


func Mp4ToMp3(buffer []byte, mimetype string) ([]byte, error) {
	if mimetype == "" {
		mimetype = "mp4"
	}
	videoFile := infra_whatsmeow_utils.GenerateTempFileName(mimetype)
	outputVideoFile := infra_whatsmeow_utils.GenerateTempFileName("mp3")

	err := os.WriteFile(videoFile, buffer, 0644)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("ffmpeg", "-i", videoFile, "-b:a", "128k", "-ac", "2", "-ar", "44100", outputVideoFile)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	buffer, err = os.ReadFile(outputVideoFile)
	if err != nil {
		return nil, err
	}

	err = os.Remove(videoFile)
	if err != nil {
		return nil, err
	}

	err = os.Remove(outputVideoFile)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
func Execute(commandProps *command_types.CommandProps) {
	quotedMsgAudio := commandProps.QuotedMsg.GetAudioMessage()
	audio := commandProps.Message.Message.GetAudioMessage()

	quotedMsgVideo := commandProps.QuotedMsg.GetVideoMessage()
	video := commandProps.Message.Message.GetVideoMessage()

	document := commandProps.QuotedMsg.GetDocumentMessage()
	quotedMsgDocument := commandProps.Message.Message.GetDocumentMessage()

	documentCaption := commandProps.QuotedMsg.GetDocumentWithCaptionMessage()
	quotedMsgDocumentCaption := commandProps.Message.Message.GetDocumentWithCaptionMessage()

	if quotedMsgVideo == nil && video == nil && document == nil && quotedMsgDocument == nil && documentCaption == nil && quotedMsgDocumentCaption == nil && quotedMsgAudio == nil && audio == nil {
		commandProps.Reply(`💬 Video para Áudio 🚀

🤔 Como usar?
✅ Marque um video com o comando!

Exemplo:
/tomp3 (mencionando)`)
		return
	}
	mediaByte := []byte{}
	if quotedMsgVideo != nil {
		if *quotedMsgVideo.FileLength > 200000000 {
			commandProps.Reply("O vídeo não pode ser maior que 200MB.")
			return
		}
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgVideo)
	} else if video != nil {
		if *video.FileLength > 200000000 {
			commandProps.Reply("O vídeo não pode ser maior que 200MB.")
			return
		}
		mediaByte, _ = commandProps.Client.Client.Download(video)
	} else if document != nil {
		if *document.FileLength > 200000000 {
			commandProps.Reply("O vídeo não pode ser maior que 200MB.")
			return
		}
		mediaByte, _ = commandProps.Client.Client.Download(document)
	} else if quotedMsgDocument != nil {
		if *quotedMsgDocument.FileLength > 200000000 {
			commandProps.Reply("O vídeo não pode ser maior que 200MB.")
			return
		}
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgDocument)
	} else if documentCaption != nil {
		if *documentCaption.Message.GetDocumentMessage().FileLength > 200000000 {
			commandProps.Reply("O vídeo não pode ser maior que 200MB.")
			return
		}
		mediaByte, _ = commandProps.Client.Client.Download(documentCaption.Message.GetDocumentMessage())
	} else if quotedMsgDocumentCaption != nil {
		if *quotedMsgDocumentCaption.Message.GetDocumentMessage().FileLength > 200000000 {
			commandProps.Reply("O vídeo não pode ser maior que 200MB.")
			return
		}
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgDocumentCaption.Message.GetDocumentMessage())
	} else if quotedMsgAudio != nil {
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgAudio)
	} else if audio != nil {
		mediaByte, _ = commandProps.Client.Client.Download(audio)
	}
	if len(mediaByte) == 0 {
		commandProps.Reply("Não foi possível converter o vídeo.")
		return
	}
	commandProps.React("👍")
	
	deviceType := infra_whatsmeow_utils.GetDevice(commandProps.Message.Info.ID)
	mimetypeByDevice := infra_whatsmeow_utils.GetMimeTypeAudioByDevice(deviceType)
	if audio != nil || quotedMsgAudio != nil {
		var ptt bool = false
		if audio != nil {
			ptt = *audio.PTT
		} else {
			ptt = *quotedMsgAudio.PTT
		}
		sender.SendAudioMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			mediaByte,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
				MimeType:      mimetypeByDevice,
				Ptt:           !ptt,
			},
		)
	} else {
		bufferMp4ToMp3, err := Mp4ToMp3(mediaByte, "mp4")
		if err != nil {
			commandProps.Reply("Não foi possível converter o vídeo.")
			return
		}
		sender.SendAudioMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			bufferMp4ToMp3,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
				MimeType:      mimetypeByDevice,
				Ptt:           false,
			},
		)
	}

}
