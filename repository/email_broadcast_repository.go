package repository

import (
	"github.com/n9mi/db-email/entity"
	"gorm.io/gorm"
)

type EmailBroadcastRepository struct {
	BaseRepository[entity.EmailBroadcast]
}

func NewEmailBroadcastRepository() *EmailBroadcastRepository {
	return new(EmailBroadcastRepository)
}

func (r *EmailBroadcastRepository) UpdateStatusByID(db *gorm.DB, id uint64, status int8) error {
	return db.Model(&entity.EmailBroadcast{}).Where("id = ?", id).Update("status", status).Error
}
