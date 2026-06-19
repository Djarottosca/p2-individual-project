package transaction

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type service struct {
	logger   *slog.Logger
	repo     Repository
	gateway  PaymentGateway
	property PropertyChecker
}

type Service interface {
	Book(userID, propertyID string, amount int64) (*Transaction, error)
	Settle(userID, bookingID string, amount int64) (*Transaction, error)
	HandlePaymentPaid(externalID string) error
	ListByUser(userID string) ([]Transaction, error)
}

func NewService(
	logger *slog.Logger,
	repo Repository,
	gateway PaymentGateway,
	property PropertyChecker,
) Service {
	return &service{logger: logger, repo: repo, gateway: gateway, property: property}
}

// Book pembeli mesan properti dengan tanda jadi.
func (s *service) Book(userID, propertyID string, amount int64) (*Transaction, error) {
	owner, status, err := s.property.GetOwnerAndStatus(propertyID)
	if err != nil {
		return nil, err
	}
	if owner == "" {
		return nil, errors.New("properti tidak ditemukan")
	}
	// tolak booking properti sendiri
	if owner == userID {
		return nil, errors.New("tidak bisa booking properti sendiri")
	}
	// tolak kalau sudah booked atau sold
	if status != "available" {
		return nil, errors.New("properti tidak tersedia untuk dibooking")
	}

	externalID := "BOOK-" + uuid.NewString()
	inv, err := s.gateway.CreateInvoice(externalID, amount, "Booking properti "+propertyID)
	if err != nil {
		s.logger.Error("gagal bikin invoice booking", "err", err)
		return nil, err
	}

	t := Transaction{
		ID:         uuid.NewString(),
		UserID:     userID,
		PropertyID: propertyID,
		Type:       "booking",
		Amount:     amount,
		Status:     "pending",
		ExternalID: externalID,
		InvoiceID:  inv.InvoiceID,
		InvoiceURL: inv.InvoiceURL,
	}
	return s.repo.Create(&t)
}

// Settle pelunasan dari booking yang sudah ada.
func (s *service) Settle(userID, bookingID string, amount int64) (*Transaction, error) {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil || booking.Type != "booking" {
		return nil, errors.New("booking tidak ditemukan")
	}
	if booking.UserID != userID {
		return nil, errors.New("bukan pemesan booking ini")
	}
	if booking.Status != "paid" {
		return nil, errors.New("booking belum dibayar, tidak bisa dilunasi")
	}

	externalID := "SETTLE-" + uuid.NewString()
	inv, err := s.gateway.CreateInvoice(externalID, amount, "Pelunasan properti "+booking.PropertyID)
	if err != nil {
		s.logger.Error("gagal bikin invoice pelunasan", "err", err)
		return nil, err
	}

	t := Transaction{
		ID:         uuid.NewString(),
		UserID:     userID,
		PropertyID: booking.PropertyID,
		BookingID:  &bookingID,
		Type:       "pelunasan",
		Amount:     amount,
		Status:     "pending",
		ExternalID: externalID,
		InvoiceID:  inv.InvoiceID,
		InvoiceURL: inv.InvoiceURL,
	}
	return s.repo.Create(&t)
}

// HandlePaymentPaid dipanggil webhook saat pembayaran lunas.
func (s *service) HandlePaymentPaid(externalID string) error {
	t, err := s.repo.FindByExternalID(externalID)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("transaksi %s tidak ditemukan", externalID)
	}
	// idempoten: kalau sudah paid, jangan diproses lagi
	if t.Status == "paid" {
		return nil
	}

	if err := s.repo.UpdateStatus(t.ID, "paid"); err != nil {
		return err
	}

	// efek menurut jenis transaksi
	switch t.Type {
	case "booking":
		return s.property.SetStatus(t.PropertyID, "booked")
	case "pelunasan":
		return s.property.SetStatus(t.PropertyID, "sold")
	}
	return nil
}

// ListByUser riwayat transaksi user.
func (s *service) ListByUser(userID string) ([]Transaction, error) {
	return s.repo.FindByUser(userID)
}

// (waktu dipakai biar import time kepakai kalau perlu nanti)
var _ = time.Now
