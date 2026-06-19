package transaction

// Invoice hasil dari payment gateway.
type Invoice struct {
	InvoiceID  string
	InvoiceURL string
}

// PaymentGateway kontrak buat bikin invoice.
// implementasinya pakai Xendit di folder repository.
type PaymentGateway interface {
	CreateInvoice(externalID string, amount int64, description string) (*Invoice, error)
}
