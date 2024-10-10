package sender

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"strings"

	infra_whatsmeow_utils "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/utils"
	"tomoribot-geminiai-version/src/utils/hooks"
	"net/http"

	"github.com/nfnt/resize"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/proto/waCommon"
	waTypes "go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func ResizeImage(image []byte, width uint, height uint) ([]byte, error) {
	imgz, err := png.Decode(strings.NewReader(string(image)))
	if err != nil {
		return nil, err
	}
	img := resize.Thumbnail(width, height, imgz, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func GetThumbVideo(video []byte) ([]byte, error) {
	videoPath := hooks.GenerateTempFileName("mp4")
	err := os.WriteFile(videoPath, video, 0644)
	if err != nil {
		return nil, err
	}
	thumbnailPath := hooks.GenerateTempFileName("png")
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01.000", "-vframes", "1", "-q:v", "2", thumbnailPath)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	thumbnail, err := os.ReadFile(thumbnailPath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(videoPath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(thumbnailPath)
	if err != nil {
		return nil, err
	}
	resized, err := ResizeImage(thumbnail, 72, 72)
	if err != nil {
		return nil, err
	}
	return resized, nil
}
func getMimeType(buffer []byte) (string, error) {
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

func SendReaction(WAClient *whatsmeow.Client, jid waTypes.JID, messageKey *waCommon.MessageKey, emoji string) (resp whatsmeow.SendResponse, err error) {
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			ReactionMessage: &waE2E.ReactionMessage{
				Key:  messageKey,
				Text: proto.String(emoji),
			},
		},
	)
}

func SendTextMessage(WAClient *whatsmeow.Client, jid waTypes.JID, body string, options *MessageOptions) (resp whatsmeow.SendResponse, mountedMessage *waE2E.Message, err error) {
	if options.QuotedMessage == nil {
		mountedMessage = &waE2E.Message{
			Conversation: proto.String(body),
		}
		sendedMsg, err := WAClient.SendMessage(
			context.Background(),
			jid.ToNonAD(),
			mountedMessage,
		)
		return sendedMsg, mountedMessage, err
	}
	user := options.QuotedMessage.Info.Sender.ToNonAD().String()
	var mentionedJid []string
	if options.ContextInfo != nil && options.ContextInfo.MentionedJID != nil {
		mentionedJid = options.ContextInfo.MentionedJID
	}
	mountedMessage = &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(body),
			ContextInfo: &waE2E.ContextInfo{
				QuotedMessage: options.QuotedMessage.Message,
				Participant:   &user,
				StanzaID:      &options.QuotedMessage.Info.ID,
				MentionedJID:  mentionedJid,
			},
		},
	}
	sendedMsg, err := WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		mountedMessage,
	)
	return sendedMsg, mountedMessage, err
}

func EditMessage(WAClient *whatsmeow.Client, jid waTypes.JID, keyId string, newText string, message *waE2E.Message) (resp whatsmeow.SendResponse, err error) {
	messageType := infra_whatsmeow_utils.GetMessageType(message)
	switch messageType {
	case "text":
		if message.ExtendedTextMessage == nil {
			message.Conversation = proto.String(newText)
		} else {
			message.ExtendedTextMessage.Text = proto.String(newText)
		}
	case "extended_text":
		message.ExtendedTextMessage.Text = proto.String(newText)
	case "image":
		message.ImageMessage.Caption = proto.String(newText)
	case "video":
		message.VideoMessage.Caption = proto.String(newText)
	}
	return WAClient.SendMessage(
		context.Background(),
		jid,
		&waE2E.Message{
			ProtocolMessage: &waE2E.ProtocolMessage{
				Key: &waCommon.MessageKey{
					FromMe:      proto.Bool(true),
					RemoteJID:   proto.String(jid.ToNonAD().String()),
					ID:          &keyId,
					Participant: proto.String(WAClient.Store.ID.ToNonAD().String()),
				},
				Type:                      waE2E.ProtocolMessage_MESSAGE_EDIT.Enum(),
				EphemeralExpiration:       proto.Uint32(0),
				EphemeralSettingTimestamp: proto.Int64(0),
				EditedMessage:             message,
				TimestampMS:               proto.Int64(0),
			},
		},
	)
}

type UploadMediaOptions struct {
	mimeType   *string
	asDocument *bool
}

func uploadMedia(client *whatsmeow.Client, buffer []byte, options *UploadMediaOptions) (uploadResponse whatsmeow.UploadResponse, mimetype string, err error) {
	typeMime := ""
	asDocument := false
	if options != nil && options.mimeType != nil && *options.mimeType != "" {
		typeMime = *options.mimeType
	} else {
		mimeType, err := getMimeType(buffer)
		if err != nil {
			fmt.Println(err)
		}
		typeMime = mimeType
	}
	if options != nil && options.asDocument != nil {
		asDocument = *options.asDocument
	}
	var appInfo whatsmeow.MediaType
	switch {
	case strings.Contains(typeMime, "image") && !asDocument:
		appInfo = whatsmeow.MediaImage
	case strings.Contains(typeMime, "video") && !strings.HasSuffix(typeMime, "webm") && !asDocument:
		appInfo = whatsmeow.MediaVideo
	case strings.Contains(typeMime, "audio") && !asDocument:
		appInfo = whatsmeow.MediaAudio
	default:
		appInfo = whatsmeow.MediaDocument
	}
	upload, err := client.Upload(context.Background(), buffer, appInfo)
	if err != nil {
		fmt.Println(err)
	}
	return upload, typeMime, nil
}

func SendImageMessage(WAClient *whatsmeow.Client, jid waTypes.JID, caption string, image []byte, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	upload, mimeType, err := uploadMedia(WAClient, image, &UploadMediaOptions{
		mimeType: &options.MimeType,
	})
	if err != nil {
		fmt.Println(err)
	}
	thumbnail, err := ResizeImage(image, 72, 72)
	if err != nil {
		thumbnail = []byte{}
	}
	if options.QuotedMessage == nil {
		return WAClient.SendMessage(
			context.Background(),
			jid.ToNonAD(),
			&waE2E.Message{
				ImageMessage: &waE2E.ImageMessage{
					Mimetype:      proto.String(mimeType),
					URL:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					JPEGThumbnail: thumbnail,
					FileEncSHA256: upload.FileEncSHA256,
					FileSHA256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
					Caption:       proto.String(caption),
				},
			},
		)
	}
	user := options.QuotedMessage.Info.Sender.ToNonAD().String()
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				Mimetype:      proto.String(mimeType),
				URL:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSHA256: upload.FileEncSHA256,
				FileSHA256:    upload.FileSHA256,
				Caption:       proto.String(caption),
				FileLength:    &upload.FileLength,
				JPEGThumbnail: thumbnail,
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaID:      &options.QuotedMessage.Info.ID,
				},
			},
		},
	)
}

func SendVideoMessage(WAClient *whatsmeow.Client, jid waTypes.JID, caption string, video []byte, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	upload, mimeType, err := uploadMedia(WAClient, video, &UploadMediaOptions{
		mimeType: &options.MimeType,
	})
	if err != nil {
		fmt.Println(err)
	}
	thumbnail, err2 := GetThumbVideo(video)
	if err2 != nil {
		thumbnail = []byte{}
	}
	if options.QuotedMessage == nil {
		return WAClient.SendMessage(
			context.Background(),
			jid.ToNonAD(),
			&waE2E.Message{
				VideoMessage: &waE2E.VideoMessage{
					Mimetype:      proto.String(mimeType),
					URL:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					FileEncSHA256: upload.FileEncSHA256,
					FileSHA256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
					JPEGThumbnail: thumbnail,
					Caption:       proto.String(caption),
					GifPlayback:   proto.Bool(options.GifPlayback),
				},
			},
		)
	}
	user := options.QuotedMessage.Info.Sender.ToNonAD().String()
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			VideoMessage: &waE2E.VideoMessage{
				Mimetype:      proto.String(mimeType),
				URL:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSHA256: upload.FileEncSHA256,
				FileSHA256:    upload.FileSHA256,
				Caption:       proto.String(caption),
				JPEGThumbnail: thumbnail,
				FileLength:    &upload.FileLength,
				GifPlayback:   proto.Bool(options.GifPlayback),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaID:      &options.QuotedMessage.Info.ID,
				},
			},
		},
	)
}

func SendAudioMessage(WAClient *whatsmeow.Client, jid waTypes.JID, audio []byte, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	mimeType := "audio/ogg"
	if options != nil && options.MimeType != "" {
		mimeType = options.MimeType
	}
	upload, mimeType, err := uploadMedia(WAClient, audio, &UploadMediaOptions{
		mimeType: &mimeType,
	})
	if err != nil {
		fmt.Println(err)
	}
	if options == nil {
		options = &MessageOptions{}
		options.Ptt = false
	}
	if options.QuotedMessage == nil {
		return WAClient.SendMessage(
			context.Background(),
			jid.ToNonAD(),
			&waE2E.Message{
				AudioMessage: &waE2E.AudioMessage{
					Mimetype:      proto.String(mimeType),
					URL:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					PTT:           proto.Bool(options.Ptt),
					MediaKey:      upload.MediaKey,
					FileEncSHA256: upload.FileEncSHA256,
					FileSHA256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
					ContextInfo:   options.ContextInfo,
				},
			},
		)
	}
	user := options.QuotedMessage.Info.Sender.ToNonAD().String()
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			AudioMessage: &waE2E.AudioMessage{
				Mimetype:      proto.String(mimeType),
				URL:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSHA256: upload.FileEncSHA256,
				PTT:           proto.Bool(options.Ptt),
				FileSHA256:    upload.FileSHA256,
				FileLength:    &upload.FileLength,
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaID:      &options.QuotedMessage.Info.ID,
				},
			},
		},
	)
}

func SendStickerMessage(WAClient *whatsmeow.Client, jid waTypes.JID, sticker []byte, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	upload, mimeType, err := uploadMedia(WAClient, sticker, &UploadMediaOptions{
		mimeType: &options.MimeType,
	})
	if err != nil {
		fmt.Println(err)
	}
	if options.QuotedMessage == nil {
		return WAClient.SendMessage(
			context.Background(),
			jid.ToNonAD(),
			&waE2E.Message{
				StickerMessage: &waE2E.StickerMessage{
					Mimetype:      proto.String(mimeType),
					URL:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					FileEncSHA256: upload.FileEncSHA256,
					FileSHA256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
					ContextInfo:   options.ContextInfo,
				},
			},
		)
	}
	user := options.QuotedMessage.Info.Sender.ToNonAD().String()
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			StickerMessage: &waE2E.StickerMessage{
				Mimetype:      proto.String(mimeType),
				URL:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSHA256: upload.FileEncSHA256,
				FileSHA256:    upload.FileSHA256,
				FileLength:    &upload.FileLength,
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaID:      &options.QuotedMessage.Info.ID,
				},
			},
		},
	)
}

func SendDocumentMessage(WAClient *whatsmeow.Client, jid waTypes.JID, FileName string, body string, document []byte, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	asDocument := true
	upload, mimeType, err := uploadMedia(WAClient, document, &UploadMediaOptions{
		mimeType:   &options.MimeType,
		asDocument: &asDocument,
	})
	if err != nil {
		fmt.Println(err)
	}
	if options.QuotedMessage == nil {
		return WAClient.SendMessage(
			context.Background(),
			jid.ToNonAD(),
			&waE2E.Message{
				DocumentMessage: &waE2E.DocumentMessage{
					Mimetype:      proto.String(mimeType),
					URL:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					FileEncSHA256: upload.FileEncSHA256,
					FileSHA256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
					ContextInfo:   options.ContextInfo,
					Title:         proto.String(options.Title),
					FileName:      proto.String(FileName),
					Caption:       proto.String(body),
				},
			},
		)
	}
	user := options.QuotedMessage.Info.Sender.ToNonAD().String()
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			DocumentMessage: &waE2E.DocumentMessage{
				Mimetype:      proto.String(mimeType),
				URL:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSHA256: upload.FileEncSHA256,
				FileSHA256:    upload.FileSHA256,
				FileLength:    &upload.FileLength,
				Title:         proto.String(options.Title),
				FileName:      proto.String(FileName),
				Caption:       proto.String(body),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaID:      &options.QuotedMessage.Info.ID,
				},
			},
		},
	)
}

func SendProtocolDeleteMessage(WAClient *whatsmeow.Client, jid waTypes.JID, user waTypes.JID, messageID string, everyone bool, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waE2E.Message{
			ProtocolMessage: &waE2E.ProtocolMessage{
				Key: &waCommon.MessageKey{
					FromMe:      proto.Bool(!everyone),
					ID:          proto.String(messageID),
					RemoteJID:   proto.String(jid.ToNonAD().String()),
					Participant: proto.String(user.ToNonAD().String()),
				},
				Type: waE2E.ProtocolMessage_REVOKE.Enum(),
			},
		},
	)
}
