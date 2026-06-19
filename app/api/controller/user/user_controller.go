package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	usersvc "p2-individual-project/service/user"
	"p2-individual-project/util/response"
)

type Controller struct {
	svc usersvc.Service
}

func NewController(svc usersvc.Service) *Controller {
	return &Controller{svc: svc}
}

type registerRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register: bind, validasi, lalu serahin ke service.
func (ctrl *Controller) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "format request tidak valid", nil)
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validasi gagal", validationErrors(err))
	}

	if err := ctrl.svc.Register(req.FullName, req.Email, req.Password); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "registrasi berhasil, cek email untuk verifikasi", nil)
}

// Verify: ambil token dari query, serahin ke service.
func (ctrl *Controller) Verify(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return response.Error(c, http.StatusBadRequest, "token wajib diisi", nil)
	}

	if err := ctrl.svc.Verify(token); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, "email terverifikasi", nil)
}

// Login: bind, validasi, balikin tokennya kalau lolos.
func (ctrl *Controller) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "format request tidak valid", nil)
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validasi gagal", validationErrors(err))
	}

	token, err := ctrl.svc.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, err.Error(), nil)
	}

	return response.Success(c, http.StatusOK, "login berhasil", map[string]string{
		"token": token,
	})
}

// validationErrors nerjemahin error validator jadi daftar field yang gagal.
func validationErrors(err error) []response.FieldError {
	var out []response.FieldError
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			out = append(out, response.FieldError{
				Field: fe.Field(),
				Issue: fe.Tag(),
			})
		}
	}
	return out
}
