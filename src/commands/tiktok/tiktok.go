package tiktok

import (
	"regexp"
	"strings"
	command_types "tomoribot-geminiai-version/src/commands/types"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"
	tiktokServices "tomoribot-geminiai-version/src/services/tiktok"
	web_functions "tomoribot-geminiai-version/src/utils/web"
)
func Details() command_types.DetailsCommand {
	return command_types.DetailsCommand{
		Name: "tiktok",
		Description: "Baixar Vídeos do TikTok sem Marca D'água",
		Category: command_types.CategoryDownload,
		Permission: command_types.PermissionAll,
		OnlyGroups: true,
		OnlyPrivate: false,
		BotRequiresAdmin: false,
		Alias: []string{"tk"},
	}
}

func CheckTikTokURL(url string) bool {
	re := regexp.MustCompile(`^(https:\/\/(www\.)?tiktok\.com\/@[\w-]+\/video\/\d+|https:\/\/vm\.tiktok\.com\/[\w-]+|https:\/\/www\.tiktok\.com\/t\/[\w-]+)`)
	return re.MatchString(url)
}

func GetTikTokURL(str string) string {
	re := regexp.MustCompile(`^(https:\/\/(www\.)?tiktok\.com\/@[\w-]+\/video\/\d+|https:\/\/vm\.tiktok\.com\/[\w-]+|https:\/\/www\.tiktok\.com\/t\/[\w-]+)`)
	return re.FindString(str)
}


func Execute(commandProps *command_types.CommandProps) {
	if !CheckTikTokURL(commandProps.Arg) {
		commandProps.Reply("Link inválido. Por favor, envie um link válido do TikTok. Exemplo: /tiktok [link]")
		return
	}
	commandProps.React("🔎")
	url := strings.Trim(GetTikTokURL(commandProps.Arg), " ")
	tiktok, err := tiktokServices.TikTokDownloader(url)
	if err != nil {
		commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
		return
	}
	var videoUrl string
	if tiktok.Play != "" {
		videoUrl = tiktok.Play
	} else if tiktok.Hdplay != "" {
		videoUrl = tiktok.Hdplay
	} else {
		commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
		return
	}
	buffer, err := web_functions.GetBufferFromUrl(videoUrl)
	if err != nil {
		commandProps.Reply("Ocorreu um erro ao baixar o vídeo. Tente novamente mais tarde.")
		return
	}
	// len than 20 mb send as document
	if len(buffer) > 20000000 {
		sender.SendDocumentMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			tiktok.Title + ".mp4",
			tiktok.Title,
			buffer,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
				MimeType: "video/mp4",
			},
		)
	} else {
		sender.SendVideoMessage(
			commandProps.Client.Client,
			commandProps.Message.Info.Chat,
			tiktok.Title,
			buffer,
			&sender.MessageOptions{
				QuotedMessage: commandProps.Message,
			},
		)
	}
}
