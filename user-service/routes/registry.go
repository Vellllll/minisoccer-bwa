package routes

import (
	"user-service/controllers"
	routes "user-service/routes/user"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.InterfaceControllerRegistry
	group      *gin.RouterGroup
}

type InterfaceRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.InterfaceControllerRegistry, group *gin.RouterGroup) InterfaceRouteRegistry {
	return &Registry{controller: controller, group: group}
}

func (r *Registry) Serve() {
	r.userRoute().Run()
}

func (r *Registry) userRoute() routes.InterfaceUserRoute {
	return routes.NewUserRoute(r.controller, r.group)
}
