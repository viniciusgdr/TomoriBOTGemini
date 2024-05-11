package instagram

import (
	"net/url"
	"regexp"
	"strings"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	pocbiServices "tomoribot-geminiai-version/src/services/pocbi"
	web_functions "tomoribot-geminiai-version/src/utils/web"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "instagram",
		Description:      "Baixar Vídeos do Instagram",
		Category:         command_types.CategoryDownload,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       true,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
		Alias:            []string{"igdl", "ig"},
	}
}

func CheckInstagramURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	if u.Host != "www.instagram.com" {
		return false
	}
	re := regexp.MustCompile(`^/reel/[^/]+/?$`)
	return re.MatchString(u.Path)
}

func Execute(commandProps *command_types.CommandProps) {
	commandProps.Arg = strings.Replace(commandProps.Arg, "/reels/", "/reel/", 1)
	if !CheckInstagramURL(commandProps.Arg) {
		commandProps.Reply("Link inválido. Por favor, envie um link válido do Instagram. Exemplo: /igdl [link]")
		return
	}
	commandProps.React("🔎")
	url := strings.Trim(commandProps.Arg, " ")
	igdl, err := pocbiServices.PocbiDownloader(url)
	if err != nil {
		commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
		return
	}
	for _, video := range igdl.Medias {
		buffer, err := web_functions.GetBufferFromUrl(video.URL)
		if err != nil {
			commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
			continue
		}
		if video.Extension == "mp4" {
			if video.Size > 25728640 {
				sender.SendDocumentMessage(
					commandProps.Client.Client,
					commandProps.Message.Info.Chat,
					igdl.Title + ".mp4",
					igdl.Title,
					buffer,
					&sender.MessageOptions{
						QuotedMessage: commandProps.Message,
						MimeType: 		"video/mp4",
					},
				)
				continue
			}
			sender.SendVideoMessage(
				commandProps.Client.Client,
				commandProps.Message.Info.Chat,
				igdl.Title,
				buffer,
				&sender.MessageOptions{
					QuotedMessage: commandProps.Message,
				},
			)
		} else {
			sender.SendImageMessage(
				commandProps.Client.Client,
				commandProps.Message.Info.Chat,
				igdl.Title,
				buffer,
				&sender.MessageOptions{
					QuotedMessage: commandProps.Message,
				},
			)
		}
	}
}
