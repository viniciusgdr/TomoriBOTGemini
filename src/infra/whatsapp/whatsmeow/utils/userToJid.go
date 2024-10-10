package infra_whatsmeow_utils

import (
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func UserToJID(user string) types.JID {
	return types.JID{
		User:   strings.Split(user, "@")[0],
		Server: strings.Split(user, "@")[1],
		Device: 0,
	}
}

func UserToJIDInternal(user string) types.JID {
	return types.JID{
		User:   strings.Split(user, "@")[0],
		Server: strings.Split(user, "@")[1],
	}
}

func UserToWhatsappJID(number string) types.JID {
	return types.JID{
		User:   number,
		Server: "s.whatsapp.net",
		Device: 0,
	}
}
