package services

import (
	"user-service/repositories"
	services "user-service/services/user"
)

type Registry struct {
	repository repositories.InterfaceRepositoryRegistry
}

type InterfaceServiceRegistry interface {
	GetUser() services.InterfaceUserService
}

func NewServiceRegistry(repository repositories.InterfaceRepositoryRegistry) InterfaceServiceRegistry {
	return &Registry{repository: repository}
}

func (r *Registry) GetUser() services.InterfaceUserService {
	return services.NewUserService(r.repository)
}
