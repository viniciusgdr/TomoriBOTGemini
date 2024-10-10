package infra_whatsmeow

import (
	"fmt"
	data_protocols "tomoribot-geminiai-version/src/data/protocols"
	constants "tomoribot-geminiai-version/src/defaults"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type whatsappConnectionRepository struct {
	client *whatsmeow.Client
	container *sqlstore.Container
}

func NewWhatsappConnectionRepository() data_protocols.DbWhatsappConnectionRepository {
	return &whatsappConnectionRepository{}
}

func (w *whatsappConnectionRepository) Start(numberPhone string) (clientReturn *whatsmeow.Client, err error) {
	container, err := sqlstore.New("sqlite3", "file:./db/database.db?_foreign_keys=on", nil)
	if err != nil {
		return nil, err
	}
	w.container = container
	deviceStores, err := container.GetAllDevices()
	var deviceStore *store.Device
	for _, device := range deviceStores {
		if device.ID.ToNonAD().User == numberPhone {
			deviceStore = device
			break
		}
	}
	if deviceStore == nil {
		deviceStore = container.NewDevice()
	}
	if err != nil {
		return nil, err
	}
	client := whatsmeow.NewClient(deviceStore, nil)
	
	if client.Store.ID == nil {
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		fmt.Println(`"⏳ Connected, but registering device ` + numberPhone + `...`)
		phoneCode, errPair := client.PairPhone(numberPhone, false, whatsmeow.PairClientChrome, constants.CLIENT_DISPLAY_NAME)
		if errPair != nil {
			return nil, errPair
		}
		fmt.Println("⏳ Waiting for code..." + phoneCode)
	} else {
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		botNumber := client.Store.ID.ToNonAD().String()
		fmt.Println("✅ Connected to WhatsApp Servers with " + botNumber)
	}

	return client, nil
}

func (w *whatsappConnectionRepository) Stop() error {
	w.client.Disconnect()
	return nil
}

func (w *whatsappConnectionRepository) Kill() error {
	err := w.client.Logout()
	return err
}

func (w *whatsappConnectionRepository) Reload() error {
	id := *w.client.Store.ID
	w.client.Disconnect()
	_, err := w.Start(id.ToNonAD().User)
	if err != nil {
		return err
	}
	return nil
}
