package sender

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"strings"

	"net/http"
	infra_whatsmeow_utils "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow/utils"
	"tomoribot-geminiai-version/src/utils/hooks"

	"github.com/nfnt/resize"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
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

func SendReaction(WAClient *whatsmeow.Client, jid waTypes.JID, messageKey *waProto.MessageKey, emoji string) (resp whatsmeow.SendResponse, err error) {
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waProto.Message{
			ReactionMessage: &waProto.ReactionMessage{
				Key:  messageKey,
				Text: proto.String(emoji),
			},
		},
	)
}

func SendInteractiveMessageV2(WAClient *whatsmeow.Client, jid waTypes.JID, title string, sections []NativeFlowListMessageSection, header *Header) (resp whatsmeow.SendResponse, mountedMessage *waProto.Message, err error) {
	msgSections := ButtonParamsJsonV2{
		Title:    "Selecione uma opção",
		Sections: sections,
	}

	mountedMessage = &waProto.Message{
		InteractiveMessage: &waProto.InteractiveMessage{
			Body: &waProto.InteractiveMessage_Body{
				Text: proto.String(title),
			},
			InteractiveMessage: &waProto.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waProto.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
					 {
						Name: proto.String("single_select"),
						ButtonParamsJson: proto.String(msgSections.toString()),
					 },
					},
					MessageParamsJson: proto.String(""),
				},
			},
		},
	}
	if header != nil && header.HasMediaAttachment {
		upload, mimeType, _ := uploadMedia(WAClient, header.MediaByte, &UploadMediaOptions{})
		thumbnail, err := ResizeImage(header.MediaByte, 72, 72)
		if err != nil {
			thumbnail = []byte{}
		}
		mountedMessage.InteractiveMessage.Header = &waProto.InteractiveMessage_Header{
			Title: proto.String(header.Title),
			Subtitle: proto.String(header.Subtitle),
			HasMediaAttachment: proto.Bool(header.HasMediaAttachment),
			Media: &waProto.InteractiveMessage_Header_ImageMessage{
				ImageMessage: &waProto.ImageMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					JpegThumbnail: thumbnail,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
				},
			},
		}
	}
	sendedMsg, err := WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		mountedMessage,
	)
	return sendedMsg, mountedMessage, err
}

func SendInteractiveMessage(WAClient *whatsmeow.Client, jid waTypes.JID, body string, buttons []InteractiveButtons, header *Header) (resp whatsmeow.SendResponse, mountedMessage *waProto.Message, err error) {
	mountedButtons := []*waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{}

	for _, button := range buttons {
		if button.Type == ButtonURL {
			nativeFlowButtonURL := NativeFlowButtonURL{
				DisplayText: button.DisplayText,
				URL:         button.ID,
				MerchantURL: button.ID,
			}
			mountedButtons = append(mountedButtons, &waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
				Name:             proto.String("cta_url"),
				ButtonParamsJson: proto.String(nativeFlowButtonURL.toString()),
			})
		} else if button.Type == ButtonCopy {
			nativeFlowButtonCopy := NativeFlowButtonCopy{
				DisplayText: button.DisplayText,
				ID:          button.ID,
				CopyCode:    button.ID,
			}
			mountedButtons = append(mountedButtons, &waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
				Name:             proto.String("cta_copy"),
				ButtonParamsJson: proto.String(nativeFlowButtonCopy.toString()),
			})
		} else if button.Type == ButtonReply {
			nativeFlowButtonReply := NativeFlowButtonReply{
				DisplayText: button.DisplayText,
				ID:          button.ID,
				Disabled:    "false",
			}
			mountedButtons = append(mountedButtons, &waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
				Name:             proto.String("quick_reply"),
				ButtonParamsJson: proto.String(nativeFlowButtonReply.toString()),
			})
		}
	}

	mountedMessage = &waProto.Message{
		InteractiveMessage: &waProto.InteractiveMessage{
			Body: &waProto.InteractiveMessage_Body{
				Text: proto.String(body),
			},
			InteractiveMessage: &waProto.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waProto.InteractiveMessage_NativeFlowMessage{
					Buttons:           mountedButtons,
					MessageParamsJson: proto.String(""),
				},
			},
		},
	}
	if header != nil && header.MediaByte != nil {
		upload, mimeType, _ := uploadMedia(WAClient, header.MediaByte, &UploadMediaOptions{})
		thumbnail, err := ResizeImage(header.MediaByte, 72, 72)
		if err != nil {
			thumbnail = []byte{}
		}
		mountedMessage.InteractiveMessage.Header = &waProto.InteractiveMessage_Header{
			Title: proto.String(header.Title),
			Subtitle: proto.String(header.Subtitle),
			HasMediaAttachment: proto.Bool(header.HasMediaAttachment),
			Media: &waProto.InteractiveMessage_Header_ImageMessage{
				ImageMessage: &waProto.ImageMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					JpegThumbnail: thumbnail,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
				},
			},
		}
	}
	sendedMsg, err := WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		mountedMessage,
	)
	return sendedMsg, mountedMessage, err
}
func SendTextMessage(WAClient *whatsmeow.Client, jid waTypes.JID, body string, options *MessageOptions) (resp whatsmeow.SendResponse, mountedMessage *waProto.Message, err error) {
	if options.QuotedMessage == nil {
		mountedMessage = &waProto.Message{
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
	if options.ContextInfo != nil && options.ContextInfo.MentionedJid != nil {
		mentionedJid = options.ContextInfo.MentionedJid
	}
	mountedMessage = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(body),
			ContextInfo: &waProto.ContextInfo{
				QuotedMessage: options.QuotedMessage.Message,
				Participant:   &user,
				StanzaId:      &options.QuotedMessage.Info.ID,
				MentionedJid:  mentionedJid,
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

func EditMessage(WAClient *whatsmeow.Client, jid waTypes.JID, keyId string, newText string, message *waProto.Message) (resp whatsmeow.SendResponse, err error) {
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
		&waProto.Message{
			ProtocolMessage: &waProto.ProtocolMessage{
				Key: &waProto.MessageKey{
					FromMe:      proto.Bool(true),
					RemoteJid:   proto.String(jid.ToNonAD().String()),
					Id:          &keyId,
					Participant: proto.String(WAClient.Store.ID.ToNonAD().String()),
				},
				Type:                      waProto.ProtocolMessage_MESSAGE_EDIT.Enum(),
				EphemeralExpiration:       proto.Uint32(0),
				EphemeralSettingTimestamp: proto.Int64(0),
				EditedMessage:             message,
				TimestampMs:               proto.Int64(0),
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
			&waProto.Message{
				ImageMessage: &waProto.ImageMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					JpegThumbnail: thumbnail,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
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
		&waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				Mimetype:      proto.String(mimeType),
				Url:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSha256: upload.FileEncSHA256,
				FileSha256:    upload.FileSHA256,
				Caption:       proto.String(caption),
				FileLength:    &upload.FileLength,
				JpegThumbnail: thumbnail,
				ContextInfo: &waProto.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaId:      &options.QuotedMessage.Info.ID,
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
			&waProto.Message{
				VideoMessage: &waProto.VideoMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
					FileLength:    &upload.FileLength,
					JpegThumbnail: thumbnail,
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
		&waProto.Message{
			VideoMessage: &waProto.VideoMessage{
				Mimetype:      proto.String(mimeType),
				Url:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSha256: upload.FileEncSHA256,
				FileSha256:    upload.FileSHA256,
				Caption:       proto.String(caption),
				JpegThumbnail: thumbnail,
				FileLength:    &upload.FileLength,
				GifPlayback:   proto.Bool(options.GifPlayback),
				ContextInfo: &waProto.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaId:      &options.QuotedMessage.Info.ID,
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
			&waProto.Message{
				AudioMessage: &waProto.AudioMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					Ptt:           proto.Bool(options.Ptt),
					MediaKey:      upload.MediaKey,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
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
		&waProto.Message{
			AudioMessage: &waProto.AudioMessage{
				Mimetype:      proto.String(mimeType),
				Url:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSha256: upload.FileEncSHA256,
				Ptt:           proto.Bool(options.Ptt),
				FileSha256:    upload.FileSHA256,
				FileLength:    &upload.FileLength,
				ContextInfo: &waProto.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaId:      &options.QuotedMessage.Info.ID,
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
			&waProto.Message{
				StickerMessage: &waProto.StickerMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
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
		&waProto.Message{
			StickerMessage: &waProto.StickerMessage{
				Mimetype:      proto.String(mimeType),
				Url:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSha256: upload.FileEncSHA256,
				FileSha256:    upload.FileSHA256,
				FileLength:    &upload.FileLength,
				ContextInfo: &waProto.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaId:      &options.QuotedMessage.Info.ID,
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
			&waProto.Message{
				DocumentMessage: &waProto.DocumentMessage{
					Mimetype:      proto.String(mimeType),
					Url:           &upload.URL,
					DirectPath:    &upload.DirectPath,
					MediaKey:      upload.MediaKey,
					FileEncSha256: upload.FileEncSHA256,
					FileSha256:    upload.FileSHA256,
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
		&waProto.Message{
			DocumentMessage: &waProto.DocumentMessage{
				Mimetype:      proto.String(mimeType),
				Url:           &upload.URL,
				DirectPath:    &upload.DirectPath,
				MediaKey:      upload.MediaKey,
				FileEncSha256: upload.FileEncSHA256,
				FileSha256:    upload.FileSHA256,
				FileLength:    &upload.FileLength,
				Title:         proto.String(options.Title),
				FileName:      proto.String(FileName),
				Caption:       proto.String(body),
				ContextInfo: &waProto.ContextInfo{
					QuotedMessage: options.QuotedMessage.Message,
					Participant:   &user,
					StanzaId:      &options.QuotedMessage.Info.ID,
				},
			},
		},
	)
}

func SendProtocolDeleteMessage(WAClient *whatsmeow.Client, jid waTypes.JID, user waTypes.JID, messageID string, everyone bool, options *MessageOptions) (resp whatsmeow.SendResponse, err error) {
	return WAClient.SendMessage(
		context.Background(),
		jid.ToNonAD(),
		&waProto.Message{
			ProtocolMessage: &waProto.ProtocolMessage{
				Key: &waProto.MessageKey{
					FromMe:      proto.Bool(!everyone),
					Id:          proto.String(messageID),
					RemoteJid:   proto.String(jid.ToNonAD().String()),
					Participant: proto.String(user.ToNonAD().String()),
				},
				Type: waProto.ProtocolMessage_REVOKE.Enum(),
			},
		},
	)
}
