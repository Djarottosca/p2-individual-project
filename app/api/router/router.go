package router

import (
	"github.com/labstack/echo/v4"

	propctrl "p2-individual-project/app/api/controller/property"
	txctrl "p2-individual-project/app/api/controller/transaction"
	userctrl "p2-individual-project/app/api/controller/user"
	mw "p2-individual-project/app/api/middleware"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	userController *userctrl.Controller,
	propertyController *propctrl.Controller,
	transactionController *txctrl.Controller,
) {
	// user, publik
	users := e.Group("/users")
	users.POST("/register", userController.Register)
	users.POST("/login", userController.Login)
	users.GET("/verify", userController.Verify)

	auth := mw.JWTAuth(jwtSecret)

	// properti
	props := e.Group("/properties")
	props.GET("", propertyController.List)
	props.GET("/:id", propertyController.Detail)
	props.POST("", propertyController.Create, auth)
	props.PUT("/:id", propertyController.Update, auth)
	props.DELETE("/:id", propertyController.Delete, auth)
	props.POST("/:id/book", transactionController.Book, auth) // booking

	e.GET("/my/properties", propertyController.ListMine, auth)

	// pelunasan
	e.POST("/bookings/:id/settle", transactionController.Settle, auth)

	// webhook, publik (dipanggil Xendit, bukan user)
	e.POST("/payments/webhook", transactionController.Webhook)

	// riwayat transaksi
	e.GET("/transactions", transactionController.List, auth)
}
