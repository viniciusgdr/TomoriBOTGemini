package domain_usecases

import "go.mau.fi/whatsmeow"

type Bot interface {
	Start(numberPhone string) (client *whatsmeow.Client, err error)
	Stop() error
	Kill() error
	Reload() error
}