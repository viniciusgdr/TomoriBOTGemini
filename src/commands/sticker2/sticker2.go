package sticker2

import (
	_ "image"
	_ "image/gif"
	"os"
	"strings"
	"tomoribot-geminiai-version/src/commands/sticker"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	"tomoribot-geminiai-version/src/services"
	"tomoribot-geminiai-version/src/utils/hooks"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "sticker2",
		Description:      "Transformar imagens ou videos em sticker em formato quadrado.",
		Category:         command_types.CategorySticker,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       true,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
		Alias:            []string{"stk2", "s2", "f2", "fig2", "figurinha2"},
	}
}

func Execute(commandProps *command_types.CommandProps) {
	image := commandProps.Message.Message.GetImageMessage()
	video := commandProps.Message.Message.GetVideoMessage()
	quotedMsgImage := commandProps.QuotedMsg.GetImageMessage()
	quotedMsgVideo := commandProps.QuotedMsg.GetVideoMessage()
	if image == nil && video == nil && quotedMsgImage == nil && quotedMsgVideo == nil {
		commandProps.Reply(`🌱 Criação de Stickers 🤖
Como fazer figurinhas?
Marque a imagem com o comando /s2

*Recursos Opcionais*:
/s2 [nome_pacote] + [nome_autor]

Use a flag --no-background para fazer um sticker transparente.
Exemplo:
/s2 --no-background`)
		return
	}
	mediaByte := []byte{}
	switch {
	case image != nil:
		mediaByte, _ = commandProps.Client.Client.Download(image)
	case video != nil:
		mediaByte, _ = commandProps.Client.Client.Download(video)
	case quotedMsgImage != nil:
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgImage)
	case quotedMsgVideo != nil:
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgVideo)
	}
	stickerPackname := "TomoriBOT WhatsApp"
	stickerAuthor := "Assine já! https://tomoribot.gdr.dev.br"
	if commandProps.Arg != "" {
		stringNoBackground := strings.Replace(commandProps.Arg, "--no-background", "", -1)
		query := strings.Split(stringNoBackground, "+")
		if len(query) > 0 {
			stickerPackname = query[0]
		}
		if len(query) > 1 {
			stickerAuthor = query[1]
		}
	}
	if image != nil || quotedMsgImage != nil {
		if strings.Contains(commandProps.Arg, "--no-background") {
			filePath := hooks.GenerateTempFileName("png")
			os.WriteFile(filePath, mediaByte, 0644)
			buffer, err := services.RemoveBg(filePath)
			if err != nil {
				commandProps.Reply("Ocorreu um erro ao remover o background.")
			}
			mediaByte = buffer
		}
		mediaWebp, _ := sticker.PngToWebp(mediaByte)

		filePath := hooks.GenerateTempFileName("webp")
		os.WriteFile(filePath, mediaWebp, 0644)

		addExif, _ := services.AddExifOnSticker(filePath, stickerPackname, stickerAuthor)
		if addExif.Success {
			mediaWebp, _ = os.ReadFile(addExif.ImagePath)
			os.Remove(addExif.ImagePath)
		}
		sender.SendStickerMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			mediaWebp,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
			},
		)
	} else if video != nil || quotedMsgVideo != nil {
		mediaWebp, _ := sticker.ToWebpGlobal(mediaByte, "mp4")

		filePath := hooks.GenerateTempFileName("webp")
		os.WriteFile(filePath, mediaWebp, 0644)

		addExif, _ := services.AddExifOnSticker(filePath, stickerPackname, stickerAuthor)
		if addExif.Success {
			mediaWebp, _ = os.ReadFile(addExif.ImagePath)
			os.Remove(addExif.ImagePath)
		}
		sender.SendStickerMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			mediaWebp,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
			},
		)
	}
}
