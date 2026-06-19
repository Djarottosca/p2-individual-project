package mailer

import "log/slog"

// LogMailer versi sementara dari Mailer.
// belum kirim email beneran, cuma nge-log token ke terminal
// biar token-nya bisa tak pake buat nembak endpoint verify pas development.
type LogMailer struct{}

func NewLogMailer() *LogMailer {
	return &LogMailer{}
}

func (m *LogMailer) SendVerification(email, token string) error {
	slog.Info("verifikasi (sementara, belum lewat email)",
		"email", email,
		"verification_token", token,
	)
	return nil
}
