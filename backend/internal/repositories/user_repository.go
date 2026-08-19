package repositories

import (
	"github.com/jj.jobo/FGC/internal/database"
	"github.com/jj.jobo/FGC/internal/models"
)

type UserRepository struct {
	*BaseRepository[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[models.User](database.DB),
	}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {

	var user models.User

	err := database.DB.
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
