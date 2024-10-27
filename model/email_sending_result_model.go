package model

import "time"

type EmailSendingRes struct {
	Recipient        string
	Error            error
	At               time.Time
	EmailBroadcastID uint64
}
