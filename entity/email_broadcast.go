package entity

import "gorm.io/gorm"

type EmailBroadcast struct {
	gorm.Model
	ID            uint64 `gorm:"primaryKey,autoIncrement"`
	EmailFormatID uint64
	EmailDest     string
	Column1Value  *string
	Column2Value  *string
	Column3Value  *string
	Column4Value  *string
	Column5Value  *string
	Status        int8
}
