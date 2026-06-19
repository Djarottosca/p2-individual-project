package user

// Mailer kontrak buat ngirim email.
// implementasinya nanti pakek Mailjet di folder repository.
type Mailer interface {
	SendVerification(email, token string) error
}
