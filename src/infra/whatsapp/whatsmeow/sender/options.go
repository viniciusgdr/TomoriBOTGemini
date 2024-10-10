package sender

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type MessageOptions struct {
	QuotedMessage *events.Message
	Footer        string
	Title         string
	MessageMedia  []byte
	ContextInfo   *waProto.ContextInfo
	GifPlayback   bool
	Ptt           bool
	MimeType      string
}
