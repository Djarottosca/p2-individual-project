package router

import (
	"github.com/labstack/echo/v4"

	userctrl "p2-individual-project/app/api/controller/user"
)

func RegisterPath(e *echo.Echo, userController *userctrl.Controller) {
	users := e.Group("/users")
	users.POST("/register", userController.Register)
	users.POST("/login", userController.Login)
	users.GET("/verify", userController.Verify)
}
