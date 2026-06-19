package main

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	propctrl "p2-individual-project/app/api/controller/property"
	txctrl "p2-individual-project/app/api/controller/transaction" // BARU
	userctrl "p2-individual-project/app/api/controller/user"
	"p2-individual-project/app/api/router"
	mailerrepo "p2-individual-project/repository/mailer"
	xenditrepo "p2-individual-project/repository/paymentGateway"
	proprepo "p2-individual-project/repository/property"
	txrepo "p2-individual-project/repository/transaction" // BARU
	userrepo "p2-individual-project/repository/user"
	propsvc "p2-individual-project/service/property"
	txsvc "p2-individual-project/service/transaction" // BARU
	usersvc "p2-individual-project/service/user"
	"p2-individual-project/util/db"
)

// CustomValidator nyambungin validator v10 ke Echo,
// biar di handler bisa panggil c.Validate(req).
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// baca .env ke environment
	if err := godotenv.Load(); err != nil {
		slog.Warn("file .env tidak ditemukan, pakai env sistem")
	}

	logger := slog.Default()

	// koneksi database
	conn := db.Connect()

	// rakit user: repo -> service -> controller
	userRepository := userrepo.NewUserRepository(conn)
	mailer := mailerrepo.NewLogMailer()
	userService := usersvc.NewService(logger, userRepository, mailer, os.Getenv("JWT_SECRET"))
	userController := userctrl.NewController(userService)

	// rakit property
	propertyRepository := proprepo.NewPropertyRepository(conn)
	propertyService := propsvc.NewService(logger, propertyRepository)
	propertyController := propctrl.NewController(propertyService)

	// rakit transaction (BARU)
	propAdapter := proprepo.NewAdapter(propertyService)
	transactionRepository := txrepo.NewTransactionRepository(conn)
	xenditGateway := xenditrepo.NewXenditGateway()
	transactionService := txsvc.NewService(logger, transactionRepository, xenditGateway, propAdapter)
	transactionController := txctrl.NewController(transactionService)

	// setup echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// daftar route, sekarang nyertain transaction controller (BARU di argumen terakhir)
	router.RegisterPath(e, os.Getenv("JWT_SECRET"), userController, propertyController, transactionController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
