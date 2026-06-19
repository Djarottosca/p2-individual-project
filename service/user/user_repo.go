package user

// Repository kontrak yang dibutuhin service buat nyentuh data user.
// implementasinya di folder repository, bukan di sini. Aku ngikutin saran Mas Pob buat merujuk yang Erdin buat.
type Repository interface {
	Create(u *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByVerificationToken(token string) (*User, error)
	MarkVerified(userID string) error
}
