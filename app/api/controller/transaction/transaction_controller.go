package transaction

import (
	"net/http"

	"github.com/labstack/echo/v4"

	txsvc "p2-individual-project/service/transaction"
	"p2-individual-project/util/response"
)

type Controller struct {
	svc txsvc.Service
}

func NewController(svc txsvc.Service) *Controller {
	return &Controller{svc: svc}
}

type bookRequest struct {
	Amount int64 `json:"amount" validate:"required,gt=0"`
}

type settleRequest struct {
	Amount int64 `json:"amount" validate:"required,gt=0"`
}

// Book: pembeli mesan properti dengan tanda jadi.
func (ctrl *Controller) Book(c echo.Context) error {
	var req bookRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "format request tidak valid", nil)
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validasi gagal", nil)
	}

	userID := c.Get("user_id").(string)
	propertyID := c.Param("id")

	tx, err := ctrl.svc.Book(userID, propertyID, req.Amount)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusCreated, "booking dibuat, silakan bayar", tx)
}

// Settle: pelunasan dari booking yang sudah dibayar.
func (ctrl *Controller) Settle(c echo.Context) error {
	var req settleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "format request tidak valid", nil)
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validasi gagal", nil)
	}

	userID := c.Get("user_id").(string)
	bookingID := c.Param("id")

	tx, err := ctrl.svc.Settle(userID, bookingID, req.Amount)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusCreated, "pelunasan dibuat, silakan bayar", tx)
}

// Webhook: dipanggil Xendit saat pembayaran lunas.
func (ctrl *Controller) Webhook(c echo.Context) error {
	// Xendit ngirim payload, kita ambil external_id dan status
	var payload struct {
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	}
	if err := c.Bind(&payload); err != nil {
		return response.Error(c, http.StatusBadRequest, "payload tidak valid", nil)
	}

	// cuma proses kalau statusnya lunas
	if payload.Status == "PAID" || payload.Status == "SETTLED" {
		if err := ctrl.svc.HandlePaymentPaid(payload.ExternalID); err != nil {
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		}
	}

	return response.Success(c, http.StatusOK, "webhook diterima", nil)
}

// List: riwayat transaksi user yang login.
func (ctrl *Controller) List(c echo.Context) error {
	userID := c.Get("user_id").(string)
	txs, err := ctrl.svc.ListByUser(userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "riwayat transaksi", txs)
}
