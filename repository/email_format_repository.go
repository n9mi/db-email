package repository

import "github.com/n9mi/db-email/entity"

type EmailFormatRepository struct {
	BaseRepository[entity.EmailFormat]
}

func NewEmailFormatRepository() *EmailFormatRepository {
	return new(EmailFormatRepository)
}
