package user

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger    *slog.Logger
	repo      Repository
	mailer    Mailer
	jwtSecret string
}

type Service interface {
	Register(fullName, email, password string) error
	Verify(token string) error
	Login(email, password string) (accessToken string, err error)
}

func NewService(
	logger *slog.Logger,
	repo Repository,
	mailer Mailer,
	jwtSecret string,
) Service {
	return &service{
		logger:    logger,
		repo:      repo,
		mailer:    mailer,
		jwtSecret: jwtSecret,
	}
}

// Register bikin akun baru, simpan, lalu kirim email verifikasi.
func (s *service) Register(fullName, email, password string) error {
	// pastiin email belum ada yang pakekk
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		s.logger.Error("gagal cek email", "err", err)
		return err
	}
	if existing != nil {
		return errors.New("email sudah terdaftar")
	}

	// hash password, jangan pernah simpan plain
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("gagal hash password", "err", err)
		return err
	}

	newUser := User{
		ID:                uuid.NewString(),
		FullName:          fullName,
		Email:             email,
		Password:          string(hashed),
		IsVerified:        false,
		VerificationToken: uuid.NewString(),
	}

	if _, err := s.repo.Create(&newUser); err != nil {
		s.logger.Error("gagal simpan user", "err", err)
		return err
	}

	// email gagal ndak ngebatalin pendaftaran, akun tetap kebuat
	if err := s.mailer.SendVerification(newUser.Email, newUser.VerificationToken); err != nil {
		s.logger.Error("gagal kirim email verifikasi", "err", err)
	}

	return nil
}

// Verify nandain akun terverifikasi lewat token dari email.
func (s *service) Verify(token string) error {
	getUser, err := s.repo.FindByVerificationToken(token)
	if err != nil {
		s.logger.Error("gagal cari token", "err", err)
		return err
	}
	if getUser == nil {
		return errors.New("token verifikasi tidak valid")
	}

	return s.repo.MarkVerified(getUser.ID)
}

// Login validasi kredensial lalu balikin JWT.
func (s *service) Login(email, password string) (string, error) {
	getUser, err := s.repo.FindByEmail(email)
	if err != nil {
		s.logger.Error("gagal cari user", "err", err)
		return "", err
	}
	if getUser == nil {
		return "", errors.New("email atau password salah")
	}

	// tolak kalau belum verifikasi
	if !getUser.IsVerified {
		return "", errors.New("email belum diverifikasi")
	}

	// nyocokin password sama hash di db
	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(password)); err != nil {
		return "", errors.New("email atau password salah")
	}

	claims := jwt.MapClaims{
		"user_id": getUser.ID,
		"email":   getUser.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("gagal bikin token", "err", err)
		return "", err
	}

	return accessToken, nil
}
