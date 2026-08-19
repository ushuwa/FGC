package repositories

import "gorm.io/gorm"

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(id uint) (*T, error) {
	var entity T

	err := r.db.First(&entity, id).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *BaseRepository[T]) FindAll(dest *[]T) error {
	return r.db.Find(dest).Error
}
