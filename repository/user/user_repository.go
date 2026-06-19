package user

import (
	"errors"
	"p2-individual-project/service/user"

	"gorm.io/gorm"
)

// UserRepository implementasi user.Repository pakai GORM.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create simpan user baru, balikin user yang udah tersimpan.
func (r *UserRepository) Create(u *user.User) (*user.User, error) {
	query := r.db.Create(u)
	if err := query.Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail cari user lewat email. kalau nggak ketemu, balikin nil tanpa error.
func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	query := r.db.Where("email = ?", email).First(&u)
	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByVerificationToken cari user lewat token verifikasi.
func (r *UserRepository) FindByVerificationToken(token string) (*user.User, error) {
	var u user.User
	query := r.db.Where("verification_token = ?", token).First(&u)
	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// MarkVerified nandain user terverifikasi, sekalian kosongin tokennya.
func (r *UserRepository) MarkVerified(userID string) error {
	query := r.db.Table("users").Where("id = ?", userID).Updates(map[string]interface{}{
		"is_verified":        true,
		"verification_token": "",
	})
	return query.Error
}
