package controllers

import (
	controllers "user-service/controllers/user"
	"user-service/services"
)

type Registry struct {
	service services.InterfaceServiceRegistry
}

type InterfaceUserController interface {
	GetUserController() controllers.InterfaceUserController
}

func NewUserController(service services.InterfaceServiceRegistry) InterfaceUserController {
	return &Registry{service: service}
}

func (u *Registry) GetUserController() controllers.InterfaceUserController {
	return controllers.NewUserController(u.service)
}
