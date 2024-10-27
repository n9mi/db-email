package entity

import "gorm.io/gorm"

type EmailFormat struct {
	gorm.Model
	ID              uint64 `gorm:"primaryKey,autoIncrement"`
	Subject         string
	BodyFormat      string
	NumCustomValue  int
	EmailBroadcasts []EmailBroadcast
}
