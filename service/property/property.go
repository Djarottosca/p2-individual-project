package property

import "time"

type Property struct {
	ID              string
	UserID          string
	Title           string
	PropertyType    string // rumah | apartemen | tanah
	TransactionType string // dijual | disewakan
	Price           int64
	LandArea        int // buat rumah & tanah
	BuildingArea    int // buat apartemen
	Bedrooms        int
	Bathrooms       int
	Certificate     string // Sertifikat Hak Milik (SHM) | Hak Guna Bangunan (HGB)
	City            string
	District        string
	Description     string
	ImageURLs       string // disimpan sebagai teks JSON
	Status          string // available | booked | sold
	FeaturedUntil   *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Property) TableName() string {
	return "properties"
}

// Filter nampung parameter pencarian dari query string.
// field yang kosong artinya nggak dipakai nyaring.
type Filter struct {
	PropertyType    string
	TransactionType string
	City            string
	District        string
	Bedrooms        int
	MinPrice        int64
	MaxPrice        int64
	MinLandArea     int
	MaxLandArea     int
	Search          string // cari di judul
	Page            int
	Limit           int
}
