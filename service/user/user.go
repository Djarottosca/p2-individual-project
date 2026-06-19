package user

import "time"

type User struct {
	ID                string
	FullName          string
	Email             string
	Password          string
	IsVerified        bool
	VerificationToken string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// TableName ngasih tahu GORM tabel mana yang dipakai struct ini.
func (User) TableName() string {
	return "users"
}
