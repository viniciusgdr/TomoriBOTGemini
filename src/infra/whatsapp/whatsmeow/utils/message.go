package infra_whatsmeow_utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	domain_models_device "tomoribot-geminiai-version/src/domain/models/device"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func GetMessageBody(v *waProto.Message) string {
	message_type := GetMessageType(v)
	var message_content string = ""

	switch message_type {
	case "text":
		if v.Conversation != nil {
			message_content = *v.Conversation
		}
	case "extended_text":
		if v.ExtendedTextMessage.Text != nil {
			message_content = *v.ExtendedTextMessage.Text
		}
	case "image":
		if v.GetViewOnceMessageV2() != nil || v.GetViewOnceMessage() != nil || v.GetViewOnceMessageV2Extension() != nil {
			message_content = ""
		} else {
			if v.ImageMessage.Caption != nil {
				message_content = *v.ImageMessage.Caption
			}
		}
	case "video":
		if v.VideoMessage.Caption != nil {
			message_content = *v.VideoMessage.Caption
		}
	case "document":
		if v.DocumentMessage.Caption != nil {
			message_content = *v.DocumentMessage.Caption
		}
	case "contact":
		if v.ContactMessage.DisplayName != nil {
			message_content = *v.ContactMessage.DisplayName
		}
	case "location":
		if v.LocationMessage.Name != nil {
			message_content = *v.LocationMessage.Name
		}
	case "buttons":
		if v.ButtonsMessage.ContentText != nil {
			message_content = *v.ButtonsMessage.ContentText
		}
	case "buttons_response":
		if v.ButtonsResponseMessage.GetSelectedDisplayText() != "" {
			message_content = v.ButtonsResponseMessage.GetSelectedDisplayText()
		}
	case "product":
		if v.ProductMessage.Body != nil {
			message_content = *v.ProductMessage.Body
		}
	case "live_location":
		if v.LiveLocationMessage.Caption != nil {
			message_content = *v.LiveLocationMessage.Caption
		}
	case "list_response":
		if v.ListResponseMessage.Title != nil {
			message_content = *v.ListResponseMessage.Title
		}
	case "template":
		templateMessage := v.TemplateMessage.GetInteractiveMessageTemplate()
		if templateMessage != nil && templateMessage.Body != nil && templateMessage.Body.Text != nil {
			message_content = *v.TemplateMessage.GetInteractiveMessageTemplate().Body.Text
		}
	case "viewOnce":
		viewOnceMessage := v.GetViewOnceMessage()
		if viewOnceMessage != nil {
			message_content = GetMessageBody(viewOnceMessage.Message)
		}
	case "viewOnceV2":
		viewOnceMessage := v.GetViewOnceMessageV2()
		if viewOnceMessage != nil {
			message_content = GetMessageBody(viewOnceMessage.Message)
		}
	case "viewOnceV2Extension":
		viewOnceMessage := v.GetViewOnceMessageV2Extension()
		if viewOnceMessage != nil {
			message_content = GetMessageBody(viewOnceMessage.Message)
		}
	case "protocol":
		protocolMessage := v.GetProtocolMessage()
		if protocolMessage != nil {
			editedMessage := protocolMessage.GetEditedMessage()
			if editedMessage != nil {
				message_content = GetMessageBody(editedMessage)
			}
		}
	case "interactiveResponseMessage":
		if v.GetInteractiveResponseMessage() != nil {
			interactiveResponseMessage := v.GetInteractiveResponseMessage()
			getInteractiveResponseMessage := interactiveResponseMessage.GetInteractiveResponseMessage().(*waProto.InteractiveResponseMessage_NativeFlowResponseMessage_)
			if getInteractiveResponseMessage != nil && getInteractiveResponseMessage.NativeFlowResponseMessage != nil && getInteractiveResponseMessage.NativeFlowResponseMessage.ParamsJson != nil {
				jsonMsg := *getInteractiveResponseMessage.NativeFlowResponseMessage.ParamsJson
				type ParamsJson struct {
					ID          string `json:"id"`
					Description string `json:"description"`
				}
				var params ParamsJson
				err := json.Unmarshal([]byte(jsonMsg), &params)
				if err != nil {
					fmt.Println(err)
				}
				message_content = params.ID
			}
		}
	}

	return message_content
}
func GetMessageType(v *waProto.Message) string {
	var result string

	switch {
	case v.GetImageMessage() != nil:
		result = "image"
	case v.GetAudioMessage() != nil:
		result = "audio"
	case v.GetVideoMessage() != nil:
		result = "video"
	case v.GetDocumentMessage() != nil:
		result = "document"
	case v.GetContactMessage() != nil:
		result = "contact"
	case v.GetLocationMessage() != nil:
		result = "location"
	case v.GetPaymentInviteMessage() != nil:
		result = "payment"
	case v.GetButtonsMessage() != nil:
		result = "buttons"
	case v.GetButtonsResponseMessage() != nil:
		result = "buttons_response"
	case v.GetProductMessage() != nil:
		result = "product"
	case v.GetLiveLocationMessage() != nil:
		result = "live_location"
	case v.GetListResponseMessage() != nil:
		result = "list_response"
	case v.GetTemplateMessage() != nil:
		result = "template"
	case v.GetExtendedTextMessage() != nil:
		result = "extended_text"
	case v.GetReactionMessage() != nil:
		result = "reaction"
	case v.GetViewOnceMessage() != nil:
		result = "viewOnce"
	case v.GetProtocolMessage() != nil:
		result = "protocol"
	case v.GetStickerMessage() != nil:
		result = "sticker"
	case v.GetConversation() != "" || v.GetExtendedTextMessage() != nil:
		result = "text"
	case v.GetViewOnceMessageV2() != nil:
		result = "viewOnceV2"
	case v.GetViewOnceMessageV2Extension() != nil:
		result = "viewOnceV2Extension"
	case v.GetInteractiveResponseMessage() != nil:
		result = "interactiveResponseMessage"
	default:
		result = "not_supported"
	}

	return result
}

func GetDevice(id string) domain_models_device.LastDevice {
	if len(id) > 22 {
		return "ANDROID"
	} else if len(id) >= 2 && id[0:2] == "3A" {
		return "IOS"
	} else {
		return "WEB"
	}
}

func GetPathTemp() string {
	return "./assets/temp/"
}
func GenerateTempFileName(ext string) string {
	ext = strings.Replace(ext, ".", "", -1)
	random := rand.Intn(1000000000)
	pathTemp := GetPathTemp()
	if _, err := os.Stat("./assets"); os.IsNotExist(err) {
		os.Mkdir("./assets", 0755)
	}
	if _, err := os.Stat(pathTemp); os.IsNotExist(err) {
		os.Mkdir(pathTemp, 0755)
	}
	if _, err := os.Stat(pathTemp + strconv.Itoa(random) + "." + ext); err == nil {
		return GenerateTempFileName(ext)
	}

	return pathTemp + strconv.Itoa(random) + "." + ext
}

func GenerateMessageID() string {
	return "TOMORIBOT" + strconv.Itoa(rand.Intn(1000000000))
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}
func GetMimeTypeAudioByDevice(device domain_models_device.LastDevice) string {
	if device == domain_models_device.ANDROID {
		return "audio/mpeg"
	} else if device == domain_models_device.IOS {
		return "audio/mp3"
	} else {
		return "audio/mpeg"
	}
}

func GetParticipant(participants []types.GroupParticipant, sender types.JID) *types.GroupParticipant {
	var groupParticipant *types.GroupParticipant
	for _, participant := range participants {
		if participant.JID == sender {
			groupParticipant = &participant
			break
		}
	}
	return groupParticipant
}

func GetAdmins(participants []types.GroupParticipant) []types.JID {
	admins := []types.JID{}
	for _, participant := range participants {
		if participant.IsAdmin {
			admins = append(admins, participant.JID)
		}
	}
	return admins
}

func GetSuperAdmin(participants []types.GroupParticipant) types.JID {
	admin := types.JID{}
	for _, participant := range participants {
		if participant.IsSuperAdmin {
			admin = participant.JID
			break
		}
	}
	return admin
}

func CheckAdminGroupNeccessary(userJid types.JID, botJid types.JID, admins []types.JID) (userChat, botChat bool) {
	adminMap := make(map[types.JID]struct{}, len(admins))
	for _, admin := range admins {
		adminMap[admin] = struct{}{}
	}
	if len(admins) < 50 {
		userChat = CheckUserIsAdminWithMap(userJid, adminMap)
		botChat = CheckUserIsAdminWithMap(botJid, adminMap)
		return userChat, botChat
	}

	resultChan := make(chan bool, 3)

	go func() {
		_, userChat = adminMap[userJid]
		resultChan <- userChat
	}()

	go func() {
		_, botChat = adminMap[botJid]
		resultChan <- botChat
	}()

	for i := 0; i < 2; i++ {
		<-resultChan
	}

	// Fechar o canal
	close(resultChan)

	return userChat, botChat
}
func CheckUserIsAdmin(jid types.JID, admins []types.JID) bool {
	adminMap := make(map[types.JID]struct{}, len(admins))
	for _, admin := range admins {
		adminMap[admin] = struct{}{}
	}
	_, ok := adminMap[jid]
	return ok
}
func CheckUserIsAdminWithMap(jid types.JID, adminMap map[types.JID]struct{}) bool {
	_, ok := adminMap[jid]
	return ok
}
func LoadMapAdmins(admins []types.JID) map[types.JID]struct{} {
	adminMap := make(map[types.JID]struct{}, len(admins))
	for _, admin := range admins {
		adminMap[admin] = struct{}{}
	}
	return adminMap
}

func GetQuotedMessage(message *waProto.Message) *waProto.Message {
	if message.ExtendedTextMessage == nil || message.ExtendedTextMessage.ContextInfo == nil || message.ExtendedTextMessage.ContextInfo.QuotedMessage == nil {
		return nil
	}
	return message.ExtendedTextMessage.ContextInfo.QuotedMessage
}

func GetQuotedMessageContextInfo(message *waProto.Message) *waProto.ContextInfo {
	if message.InteractiveResponseMessage != nil && message.InteractiveResponseMessage.GetInteractiveResponseMessage() != nil {
		return message.InteractiveResponseMessage.ContextInfo
	}
	quotedMessage := GetQuotedMessage(message)
	if quotedMessage == nil {
		return nil
	}
	return message.ExtendedTextMessage.ContextInfo
}
