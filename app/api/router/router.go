package router

import (
	"github.com/labstack/echo/v4"

	propctrl "p2-individual-project/app/api/controller/property"
	userctrl "p2-individual-project/app/api/controller/user"
	mw "p2-individual-project/app/api/middleware"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	userController *userctrl.Controller,
	propertyController *propctrl.Controller,
) {
	// user, semua publik
	users := e.Group("/users")
	users.POST("/register", userController.Register)
	users.POST("/login", userController.Login)
	users.GET("/verify", userController.Verify)

	// middleware buat endpoint yang butuh login
	auth := mw.JWTAuth(jwtSecret)

	// properti
	props := e.Group("/properties")
	props.GET("", propertyController.List)                // publik
	props.GET("/:id", propertyController.Detail)          // publik
	props.POST("", propertyController.Create, auth)       // butuh login
	props.PUT("/:id", propertyController.Update, auth)    // butuh login
	props.DELETE("/:id", propertyController.Delete, auth) // butuh login

	// listing milik sendiri
	e.GET("/my/properties", propertyController.ListMine, auth)
}
