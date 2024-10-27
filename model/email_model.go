package model

type EmailModel struct {
	To               string
	Subject          string
	Body             string
	EmailBroadcastID uint64
}
