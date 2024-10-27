package repository

import "gorm.io/gorm"

type BaseRepository[T any] struct {
}

func (r *BaseRepository[T]) FindByID(db *gorm.DB, e *T, ID any) error {
	return db.First(e, "id = ?", ID).Error
}

func (r *BaseRepository[T]) IsExistsByID(db *gorm.DB, ID any) (bool, error) {
	var count int
	err := db.Model(new(T)).
		Select("count(1)").
		Where("id = ?", ID).
		Find(&count).
		Error

	return count == 1, err
}

func (r *BaseRepository[T]) Create(db *gorm.DB, e *T) error {
	return db.Create(e).Error
}

func (r *BaseRepository[T]) CreateMany(db *gorm.DB, e []*T) error {
	return db.Create(e).Error
}

func (r *BaseRepository[T]) Update(db *gorm.DB, e *T, ID any) error {
	return db.Model(e).Updates(e).Error
}

func (r *BaseRepository[T]) Delete(db *gorm.DB, e *T) error {
	return db.Delete(e).Error
}
