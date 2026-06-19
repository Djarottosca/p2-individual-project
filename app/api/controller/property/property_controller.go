package property

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	propsvc "p2-individual-project/service/property"
	"p2-individual-project/util/response"
)

type Controller struct {
	svc propsvc.Service
}

func NewController(svc propsvc.Service) *Controller {
	return &Controller{svc: svc}
}

// request body buat create & update.
type propertyRequest struct {
	Title           string `json:"title" validate:"required"`
	PropertyType    string `json:"property_type" validate:"required,oneof=rumah apartemen tanah"`
	TransactionType string `json:"transaction_type" validate:"required,oneof=dijual disewakan"`
	Price           int64  `json:"price" validate:"required,gt=0"`
	LandArea        int    `json:"land_area"`
	BuildingArea    int    `json:"building_area"`
	Bedrooms        int    `json:"bedrooms"`
	Bathrooms       int    `json:"bathrooms"`
	Certificate     string `json:"certificate"`
	City            string `json:"city" validate:"required"`
	District        string `json:"district" validate:"required"`
	Description     string `json:"description"`
	ImageURLs       string `json:"image_urls"`
}

// Create: pasang listing. pemilik diambil dari token, bukan body.
func (ctrl *Controller) Create(c echo.Context) error {
	var req propertyRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "format request tidak valid", nil)
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validasi gagal", validationErrors(err))
	}

	userID := c.Get("user_id").(string)

	created, err := ctrl.svc.Create(userID, toDomain(req))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusCreated, "listing berhasil dibuat", created)
}

// List: daftar properti dengan filter dari query string.
func (ctrl *Controller) List(c echo.Context) error {
	f := propsvc.Filter{
		PropertyType:    c.QueryParam("property_type"),
		TransactionType: c.QueryParam("transaction_type"),
		City:            c.QueryParam("city"),
		District:        c.QueryParam("district"),
		Search:          c.QueryParam("q"),
		Bedrooms:        atoi(c.QueryParam("bedrooms")),
		MinPrice:        atoi64(c.QueryParam("min_price")),
		MaxPrice:        atoi64(c.QueryParam("max_price")),
		MinLandArea:     atoi(c.QueryParam("min_land_area")),
		MaxLandArea:     atoi(c.QueryParam("max_land_area")),
		Page:            atoi(c.QueryParam("page")),
		Limit:           atoi(c.QueryParam("limit")),
	}

	props, err := ctrl.svc.List(f)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "daftar properti", props)
}

// Detail: satu listing.
func (ctrl *Controller) Detail(c echo.Context) error {
	id := c.Param("id")
	p, err := ctrl.svc.Detail(id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "detail properti", p)
}

// Update: ubah listing milik sendiri.
func (ctrl *Controller) Update(c echo.Context) error {
	var req propertyRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "format request tidak valid", nil)
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validasi gagal", validationErrors(err))
	}

	userID := c.Get("user_id").(string)
	id := c.Param("id")

	updated, err := ctrl.svc.Update(userID, id, toDomain(req))
	if err != nil {
		// bedain "bukan pemilik" jadi 403
		if err.Error() == "bukan pemilik listing" {
			return response.Error(c, http.StatusForbidden, err.Error(), nil)
		}
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "listing berhasil diubah", updated)
}

// Delete: hapus listing milik sendiri.
func (ctrl *Controller) Delete(c echo.Context) error {
	userID := c.Get("user_id").(string)
	id := c.Param("id")

	if err := ctrl.svc.Delete(userID, id); err != nil {
		if err.Error() == "bukan pemilik listing" {
			return response.Error(c, http.StatusForbidden, err.Error(), nil)
		}
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "listing berhasil dihapus", nil)
}

// ListMine: listing milik user yang lagi login.
func (ctrl *Controller) ListMine(c echo.Context) error {
	userID := c.Get("user_id").(string)
	props, err := ctrl.svc.ListMine(userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "listing milik saya", props)
}

// toDomain ngubah request jadi struct domain.
func toDomain(req propertyRequest) propsvc.Property {
	return propsvc.Property{
		Title:           req.Title,
		PropertyType:    req.PropertyType,
		TransactionType: req.TransactionType,
		Price:           req.Price,
		LandArea:        req.LandArea,
		BuildingArea:    req.BuildingArea,
		Bedrooms:        req.Bedrooms,
		Bathrooms:       req.Bathrooms,
		Certificate:     req.Certificate,
		City:            req.City,
		District:        req.District,
		Description:     req.Description,
		ImageURLs:       req.ImageURLs,
	}
}

// helper kecil buat ngubah query string ke angka, kosong jadi 0.
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func atoi64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

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
