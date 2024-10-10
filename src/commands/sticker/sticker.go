package sticker

import (
	"bytes"
	_ "image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	"tomoribot-geminiai-version/src/services"
	"tomoribot-geminiai-version/src/utils/hooks"

	"github.com/chai2010/webp"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "sticker",
		Description:      "Transformar imagens ou videos em sticker.",
		Category:         command_types.CategorySticker,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       true,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
		Alias:            []string{"stk", "s", "f", "fig", "figurinha"},
	}
}

func PngToWebp(pngData []byte) ([]byte, error) {
	image, err := jpeg.Decode(bytes.NewReader(pngData))
	if err != nil {
		pngImage, errPng := png.Decode(bytes.NewReader(pngData))
		if errPng != nil {
			return nil, errPng
		}
		image = pngImage
	}
	var webpBuf bytes.Buffer

	err = webp.Encode(&webpBuf, image, nil)
	if err != nil {
		return nil, err
	}
	return webpBuf.Bytes(), nil
}
func Random(start int, end int) int {
	return rand.Intn(end-start) + start
}
func ToWebpGlobal(videoData []byte, mimetypeInput string) ([]byte, error) {
	inputFilePath := hooks.GenerateTempFileName(mimetypeInput)
	outputFilePath := hooks.GenerateTempFileName("webp")
	err := os.WriteFile(inputFilePath, videoData, 0644)
	if err != nil {
		return nil, err
	}
	err = exec.Command(
		"ffmpeg", "-i",
		inputFilePath,
		"-vcodec",
		"libwebp",
		"-vf",
		"scale=\\'iw*min(300/iw\\,300/ih)\\':\\'ih*min(300/iw\\,300/ih)\\',format=rgba,pad=300:300:\\'(300-iw)/2\\':\\'(300-ih)/2\\':\\'#00000000\\',setsar=1,fps=15",
		"-loop",
		"0",
		"-ss",
		"00:00:00.0",
		"-t",
		"00:00:05.0",
		"-preset",
		"default",
		"-an",
		"-vsync",
		"0",
		"-s",
		"512:512",
		outputFilePath).Run()

	if err != nil {
		return nil, err
	}
	webpFile, err := os.ReadFile(outputFilePath)
	if err != nil {
		return nil, err
	}
	os.Remove(inputFilePath)
	os.Remove(outputFilePath)
	return webpFile, nil
}

func Execute(commandProps *command_types.CommandProps) {
	image := commandProps.Message.Message.GetImageMessage()
	video := commandProps.Message.Message.GetVideoMessage()
	quotedMsgImage := commandProps.QuotedMsg.GetImageMessage()
	quotedMsgVideo := commandProps.QuotedMsg.GetVideoMessage()
	quotedMsgSticker := commandProps.QuotedMsg.GetStickerMessage()
	if image == nil && video == nil && quotedMsgImage == nil && quotedMsgVideo == nil && quotedMsgSticker == nil {
		commandProps.Reply(`🌱 Criação de Stickers 🤖
Como fazer figurinhas?
Marque a imagem com o comando /s

*Recursos Opcionais*:
/s [nome_pacote] + [nome_autor]

Use a flag --no-background para fazer um sticker transparente.
Exemplo:
/s --no-background`)
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
	case quotedMsgSticker != nil:
		mediaByte, _ = commandProps.Client.Client.Download(quotedMsgSticker)
	}
	stickerPackname := "TomoriBOT WhatsApp"
	stickerAuthor := "Feito usando o TomoriBOT"
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
	if image != nil || quotedMsgImage != nil || quotedMsgSticker != nil {
		if strings.Contains(commandProps.Arg, "--no-background") {
			filePath := hooks.GenerateTempFileName("png")
			os.WriteFile(filePath, mediaByte, 0644)
			buffer, err := services.RemoveBg(filePath)
			if err != nil {
				commandProps.Reply("Ocorreu um erro ao remover o background, pulando para o próximo passo.")
			}
			mediaByte = buffer
		}
		mediaWebp, _ := ToWebpGlobal(mediaByte, "png")

		filePath := hooks.GenerateTempFileName("webp")
		os.WriteFile(filePath, mediaWebp, 0644)

		addExif, err := services.AddExifOnSticker(filePath, stickerPackname, stickerAuthor)
		if err != nil {
			os.Remove(filePath)
		} else if addExif.Success {
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
		mediaWebp, _ := ToWebpGlobal(mediaByte, "mp4")

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
