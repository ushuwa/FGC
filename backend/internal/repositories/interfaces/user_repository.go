package interfaces

import "github.com/jj.jobo/FGC/internal/models"

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
}
