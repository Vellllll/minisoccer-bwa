package controllers

import (
	controllers "user-service/controllers/user"
	"user-service/services"
)

type Registry struct {
	service services.InterfaceServiceRegistry
}

type InterfaceControllerRegistry interface {
	GetUserController() controllers.InterfaceUserController
}

func NewUserController(service services.InterfaceServiceRegistry) InterfaceControllerRegistry {
	return &Registry{service: service}
}

func (u *Registry) GetUserController() controllers.InterfaceUserController {
	return controllers.NewUserController(u.service)
}
