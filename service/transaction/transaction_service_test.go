package transaction

import (
	"errors"
	"log/slog"
	"testing"
)

// helper bikin service dengan tiga mock.
func newTestService(repo Repository, gw PaymentGateway, prop PropertyChecker) Service {
	return NewService(slog.Default(), repo, gw, prop)
}

func TestBook_PropertiSendiri(t *testing.T) {
	// properti dimiliki user yang sama yang mau booking -> harus ditolak
	prop := &mockProperty{
		getOwnerAndStatusFunc: func(propertyID string) (string, string, error) {
			return "user-1", "available", nil
		},
	}
	repo := &mockRepo{}
	gw := &mockGateway{}

	svc := newTestService(repo, gw, prop)

	_, err := svc.Book("user-1", "prop-1", 1000000)
	if err == nil {
		t.Fatal("harusnya error karena booking properti sendiri, tapi nil")
	}
}

func TestBook_PropertiTidakTersedia(t *testing.T) {
	// properti udah booked -> harus ditolak
	prop := &mockProperty{
		getOwnerAndStatusFunc: func(propertyID string) (string, string, error) {
			return "user-2", "booked", nil
		},
	}
	repo := &mockRepo{}
	gw := &mockGateway{}

	svc := newTestService(repo, gw, prop)

	_, err := svc.Book("user-1", "prop-1", 1000000)
	if err == nil {
		t.Fatal("harusnya error karena properti tidak tersedia, tapi nil")
	}
}

func TestBook_Sukses(t *testing.T) {
	// properti milik orang lain & available -> booking sukses
	prop := &mockProperty{
		getOwnerAndStatusFunc: func(propertyID string) (string, string, error) {
			return "user-2", "available", nil
		},
	}
	gw := &mockGateway{
		createInvoiceFunc: func(externalID string, amount int64, description string) (*Invoice, error) {
			return &Invoice{InvoiceID: "inv-1", InvoiceURL: "https://xendit/inv-1"}, nil
		},
	}
	var created bool
	repo := &mockRepo{
		createFunc: func(tr *Transaction) (*Transaction, error) {
			created = true
			return tr, nil
		},
	}

	svc := newTestService(repo, gw, prop)

	tr, err := svc.Book("user-1", "prop-1", 1000000)
	if err != nil {
		t.Fatalf("harusnya sukses, tapi error: %v", err)
	}
	if !created {
		t.Fatal("harusnya transaksi dibuat, tapi Create nggak kepanggil")
	}
	if tr.Status != "pending" {
		t.Fatalf("status harusnya pending, dapat: %s", tr.Status)
	}
	if tr.InvoiceURL == "" {
		t.Fatal("invoice url harusnya keisi")
	}
}

// biar import errors kepakai kalau nanti butuh
var _ = errors.New
