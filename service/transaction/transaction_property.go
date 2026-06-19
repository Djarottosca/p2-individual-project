package transaction

// PropertyReader kontrak minimal buat ngecek & ngubah properti.
// service property nanti yang implementasiin lewat adapter di main.
type PropertyChecker interface {
	GetOwnerAndStatus(propertyID string) (ownerID string, status string, err error)
	SetStatus(propertyID, status string) error
}
