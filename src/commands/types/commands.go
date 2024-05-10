package command_types

import (
	"time"
	"tomoribot-geminiai-version/client"
	"tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/sender"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type Command struct {
	Details DetailsCommand
	Execute func(commandProps *CommandProps)
}

type Permission string

const (
	PermissionAdmin Permission = "admin"
	PermissionOwner Permission = "private"
	PermissionAll   Permission = "all"
)

type Category string

const (
	CategoryGroup    Category = "group"
	CategoryDownload Category = "download"
	CategorySticker  Category = "sticker"
	CategoryOffTopic Category = "off-topic"
	CategoryUtility  Category = "utility"
	CategoryInfo     Category = "info"
	CategorySearch   Category = "search"
	CategoryFun      Category = "fun"
)

type DetailsCommand struct {
	Name             string
	Description      string
	Category         Category
	Permission       Permission
	OnlyGroups       bool
	OnlyPrivate      bool
	BotRequiresAdmin bool
	Alias            []string
}

type CommandProps struct {
	Client               *client.Client
	Args                 []string
	Message              *events.Message
	QuotedMsgContextInfo *waProto.ContextInfo
	Arg                  string
	Timestamp            time.Time
	UserChat             bool
	BotChat              bool
	MessageType          string
}

func (props *CommandProps) Reply(text string) (whatsmeow.SendResponse, *waProto.Message, error) {
	options := &sender.MessageOptions{
		QuotedMessage: props.Message,
	}
	return sender.SendTextMessage(
		props.Client.Client,
		props.Message.Info.Chat,
		text,
		options,
	)
}

func (props *CommandProps) React(emoji string) (whatsmeow.SendResponse, error) {
	participant := props.Message.Info.Sender.ToNonAD().String()
	groupId := props.Message.Info.Chat.ToNonAD().String()
	messagekey := &waProto.MessageKey{
		Id:        &props.Message.Info.ID,
		FromMe:    proto.Bool(false),
		RemoteJid: &groupId,
	}
	if props.Message.Info.IsGroup {
		messagekey.Participant = &participant
	}
	if props.Message.Info.IsFromMe {
		messagekey.FromMe = proto.Bool(true)
	}
	return sender.SendReaction(
		props.Client.Client,
		props.Message.Info.Chat,
		messagekey,
		emoji,
	)
}
