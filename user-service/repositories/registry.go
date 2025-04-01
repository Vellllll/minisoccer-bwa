package repositories

import (
	repositories "user-service/repositories/user"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type InterfaceRepositoryRegistry interface {
	GetUser() repositories.InterfaceUserRepository
}

func NewRepositoryRegistry(db *gorm.DB) InterfaceRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetUser() repositories.InterfaceUserRepository {
	return repositories.NewUserRepository(r.db)
}
