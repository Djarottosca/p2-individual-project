package transaction

type Repository interface {
	Create(t *Transaction) (*Transaction, error)
	FindByExternalID(externalID string) (*Transaction, error)
	FindByID(id string) (*Transaction, error)
	UpdateStatus(id, status string) error
	FindByUser(userID string) ([]Transaction, error)
}
