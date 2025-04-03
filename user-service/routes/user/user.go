package routes

import (
	"user-service/controllers"
	"user-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.InterfaceControllerRegistry
	group      *gin.RouterGroup
}

type InterfaceUserRoute interface {
	Run()
}

func NewUserRoute(controller controllers.InterfaceControllerRegistry, group *gin.RouterGroup) InterfaceUserRoute {
	return &UserRoute{controller: controller, group: group}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.GET("/user", middlewares.Authenticate(), u.controller.GetUserController().GetUserLogin)
	group.GET("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().GetUserByUUID)
	group.POST("/login", u.controller.GetUserController().Login)
	group.POST("/register", u.controller.GetUserController().Register)
	group.PUT("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().Update)
}
