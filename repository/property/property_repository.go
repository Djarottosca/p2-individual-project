package property // implementasi GORM

import (
	"errors"

	"p2-individual-project/service/property"

	"gorm.io/gorm"
)

type PropertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

func (r *PropertyRepository) Create(p *property.Property) (*property.Property, error) {
	if err := r.db.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

// FindAll nyusun query bertahap sesuai filter yang keisi.
func (r *PropertyRepository) FindAll(f property.Filter) ([]property.Property, error) {
	var props []property.Property
	query := r.db.Model(&property.Property{})

	// filter pilih-dari-daftar: cuma dipasang kalau keisi
	if f.PropertyType != "" {
		query = query.Where("property_type = ?", f.PropertyType)
	}
	if f.TransactionType != "" {
		query = query.Where("transaction_type = ?", f.TransactionType)
	}
	if f.City != "" {
		query = query.Where("city = ?", f.City)
	}
	if f.District != "" {
		query = query.Where("district = ?", f.District)
	}
	if f.Bedrooms > 0 {
		query = query.Where("bedrooms >= ?", f.Bedrooms)
	}

	// filter rentang
	if f.MinPrice > 0 {
		query = query.Where("price >= ?", f.MinPrice)
	}
	if f.MaxPrice > 0 {
		query = query.Where("price <= ?", f.MaxPrice)
	}
	if f.MinLandArea > 0 {
		query = query.Where("land_area >= ?", f.MinLandArea)
	}
	if f.MaxLandArea > 0 {
		query = query.Where("land_area <= ?", f.MaxLandArea)
	}

	// pencarian judul
	if f.Search != "" {
		query = query.Where("title ILIKE ?", "%"+f.Search+"%")
	}

	// featured yang masih aktif naik ke atas, lalu yang terbaru
	query = query.Order("featured_until DESC NULLS LAST").Order("created_at DESC")

	// pagination
	offset := (f.Page - 1) * f.Limit
	query = query.Limit(f.Limit).Offset(offset)

	if err := query.Find(&props).Error; err != nil {
		return nil, err
	}
	return props, nil
}

func (r *PropertyRepository) FindByID(id string) (*property.Property, error) {
	var p property.Property
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *PropertyRepository) Update(p *property.Property) (*property.Property, error) {
	if err := r.db.Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PropertyRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&property.Property{}).Error
}

func (r *PropertyRepository) FindByOwner(userID string) ([]property.Property, error) {
	var props []property.Property
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&props).Error
	if err != nil {
		return nil, err
	}
	return props, nil
}
