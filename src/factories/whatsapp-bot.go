package factories

import (
	data_usecases "tomoribot-geminiai-version/src/data/usecases"
	domain_usecases "tomoribot-geminiai-version/src/domain/usecases"
	infra_whatsmeow "tomoribot-geminiai-version/src/infra/whatsapp/whatsmeow"
)

func MakeWhatsappBot() domain_usecases.Bot {
	whatsappConnectionRepository := infra_whatsmeow.NewWhatsappConnectionRepository()
	dbWhatsAppBot := data_usecases.NewDbWhatsappBot(whatsappConnectionRepository)
	return dbWhatsAppBot
}