package transaction

import "time"

type Transaction struct {
	ID         string
	UserID     string
	PropertyID string
	BookingID  *string // diisi pas pelunasan, nunjuk booking-nya
	Type       string  // booking | pelunasan
	Amount     int64
	Status     string // pending | paid | expired
	ExternalID string // id yang kita kirim ke xendit
	InvoiceID  string // id balikan dari xendit
	InvoiceURL string // link pembayaran
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Transaction) TableName() string {
	return "transactions"
}
