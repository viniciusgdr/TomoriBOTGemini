package client

import (
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"gorm.io/gorm"
)

type Client struct {
	DB            *gorm.DB
	DBGroupMutex  *sync.Mutex
	DBMemberMutex *sync.Mutex
	Client        *whatsmeow.Client
	StartedAt     time.Time
}
