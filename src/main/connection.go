package loader

import (
	"fmt"
	"sync"
	"time"
	"tomoribot-geminiai-version/client"
	"tomoribot-geminiai-version/src/factories"
	"tomoribot-geminiai-version/src/handlers"

	"go.mau.fi/whatsmeow/types/events"
	"gorm.io/gorm"
)

func StartConnection(
	phoneNumber string,
	DB *gorm.DB,
) (clientReturn *client.Client) {
	whatsAppBot := factories.MakeWhatsappBot()
	clientWhatsMeow, err := whatsAppBot.Start(phoneNumber)
	if err != nil {
		time.Sleep(5 * time.Second)
		return StartConnection(phoneNumber, DB)
	}
	clientType := &client.Client{
		DB:            DB,
		Client:        clientWhatsMeow,
		StartedAt:     time.Now(),
		DBGroupMutex:  &sync.Mutex{},
		DBMemberMutex: &sync.Mutex{},
	}

	clientWhatsMeow.AddEventHandler(func(evt interface{}) {
		switch e := evt.(type) {
		case *events.Connected:
			{
				fmt.Println("✅ Connected to WhatsApp Servers and authenticated!")
			}
		case *events.Disconnected:
			{
				fmt.Println("❌ Disconnected from WhatsApp Servers!")
				time.Sleep(5 * time.Second)
				StartConnection(phoneNumber, DB)
			}
		case *events.Message:
			{
				go handlers.MessageHandler(clientType, e)
			}
		}
	})

	return clientType
}
