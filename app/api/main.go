package main

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	propctrl "p2-individual-project/app/api/controller/property"
	userctrl "p2-individual-project/app/api/controller/user"
	"p2-individual-project/app/api/router"
	mailerrepo "p2-individual-project/repository/mailer"
	proprepo "p2-individual-project/repository/property"
	userrepo "p2-individual-project/repository/user"
	propsvc "p2-individual-project/service/property"
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

	// rakit user: repo -> service -> controller, service tidak depend ke repo, define contract repository di service layer
	userRepository := userrepo.NewUserRepository(conn)
	mailer := mailerrepo.NewLogMailer()
	userService := usersvc.NewService(logger, userRepository, mailer, os.Getenv("JWT_SECRET"))
	userController := userctrl.NewController(userService)

	// ngerakit lapisan property
	propertyRepository := proprepo.NewPropertyRepository(conn)
	propertyService := propsvc.NewService(logger, propertyRepository)
	propertyController := propctrl.NewController(propertyService)

	// setup echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// daftar route
	router.RegisterPath(e, os.Getenv("JWT_SECRET"), userController, propertyController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
