package twitter

import (
	"regexp"
	"strings"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	twitterServices "tomoribot-geminiai-version/src/services/twitter"
	web_functions "tomoribot-geminiai-version/src/utils/web"
)

func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name:             "twitter",
		Description:      "Baixar Vídeos do Twitter em HD",
		Category:         command_types.CategoryDownload,
		Permission:       command_types.PermissionAll,
		OnlyGroups:       true,
		OnlyPrivate:      false,
		BotRequiresAdmin: false,
	Alias:            []string{"tw", "x"},
	}
}

func CheckTwitterURL(url string) bool {
	re := regexp.MustCompile(`^https?://twitter.com/[a-zA-Z0-9_]{1,15}/status/\d+(\?.*)?$`)
	reX := regexp.MustCompile(`^https?://x.com/[a-zA-Z0-9_]{1,15}/status/\d+(\?.*)?$`)

	return re.MatchString(url) || reX.MatchString(url)
}

func Execute(commandProps *command_types.CommandProps) {
	if !CheckTwitterURL(commandProps.Arg) {
		commandProps.Reply("Link inválido. Por favor, envie um link válido do Twitter. Exemplo: /twitter [link]")
		return
	}
	commandProps.React("🔎")
	url := strings.Trim(commandProps.Arg, " ")
	url = strings.Replace(url, "x.com", "twitter.com", 1)
	twitter, err := twitterServices.TwitterOfficial(url)
	if err != nil {
		twitter2, err2 := twitterServices.TwitterDownloader(url)
		if err2 != nil {
			commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
			return
		}
		twitter = twitter2
	}
	buffer, err := web_functions.GetBufferFromUrl(twitter.Url)
	if err != nil {
		commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
		return
	}
	commandProps.React("✅")
	if twitter.Quality == "image" {
		sender.SendImageMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			"",
			buffer,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
			},
		)
	} else {
		sender.SendVideoMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			"",
			buffer,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
				MimeType:	"video/mp4",
			},
		)
	}
}
