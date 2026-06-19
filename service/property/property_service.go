package property

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

// ini business logic

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	Create(userID string, p Property) (*Property, error)
	List(f Filter) ([]Property, error)
	Detail(id string) (*Property, error)
	Update(userID, id string, p Property) (*Property, error)
	Delete(userID, id string) error
	ListMine(userID string) ([]Property, error)
	SetStatus(id, status string) error
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

// Create pasang listing baru. pemilik diambil dari userID (token), bukan dari input.
func (s *service) Create(userID string, p Property) (*Property, error) {
	p.ID = uuid.NewString()
	p.UserID = userID
	p.Status = "available"

	created, err := s.repo.Create(&p)
	if err != nil {
		s.logger.Error("gagal simpan properti", "err", err)
		return nil, err
	}
	return created, nil
}

// List ngembaliin daftar properti sesuai filter.
func (s *service) List(f Filter) ([]Property, error) {
	// kasih default pagination biar nggak narik semua sekaligus
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 10
	}
	return s.repo.FindAll(f)
}

// Detail ngambil satu properti.
func (s *service) Detail(id string) (*Property, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("gagal ambil properti", "err", err)
		return nil, err
	}
	if p == nil {
		return nil, errors.New("properti tidak ditemukan")
	}
	return p, nil
}

// Update ubah listing. cuma pemilik yang boleh.
func (s *service) Update(userID, id string, input Property) (*Property, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("properti tidak ditemukan")
	}

	// ownership check: tolak kalau bukan punya dia
	if existing.UserID != userID {
		return nil, errors.New("bukan pemilik listing")
	}

	// timpa field yang boleh diubah, jangan sentuh id, pemilik, status, featured
	existing.Title = input.Title
	existing.PropertyType = input.PropertyType
	existing.TransactionType = input.TransactionType
	existing.Price = input.Price
	existing.LandArea = input.LandArea
	existing.BuildingArea = input.BuildingArea
	existing.Bedrooms = input.Bedrooms
	existing.Bathrooms = input.Bathrooms
	existing.Certificate = input.Certificate
	existing.City = input.City
	existing.District = input.District
	existing.Description = input.Description
	existing.ImageURLs = input.ImageURLs

	updated, err := s.repo.Update(existing)
	if err != nil {
		s.logger.Error("gagal update properti", "err", err)
		return nil, err
	}
	return updated, nil
}

// Delete hapus listing. cuma pemilik yang boleh.
func (s *service) Delete(userID, id string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("properti tidak ditemukan")
	}
	if existing.UserID != userID {
		return errors.New("bukan pemilik listing")
	}
	return s.repo.Delete(id)
}

// ListMine ngembaliin listing milik user yang lagi login.
func (s *service) ListMine(userID string) ([]Property, error) {
	return s.repo.FindByOwner(userID)
}

// SetStatus ngubah status properti, dipakai alur pembayaran.
func (s *service) SetStatus(id, status string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("properti tidak ditemukan")
	}
	existing.Status = status
	_, err = s.repo.Update(existing)
	return err
}
