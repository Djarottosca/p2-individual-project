package transaction

// File ini mock manual buat keperluan unit test usecase.
// nggak pakai generator, cukup struct yang memenuhi tiap interface,
// lalu kita atur balikannya pakai field fungsi.

// --- mock Repository ---

type mockRepo struct {
	createFunc           func(t *Transaction) (*Transaction, error)
	findByExternalIDFunc func(externalID string) (*Transaction, error)
	findByIDFunc         func(id string) (*Transaction, error)
	updateStatusFunc     func(id, status string) error
	findByUserFunc       func(userID string) ([]Transaction, error)
}

func (m *mockRepo) Create(t *Transaction) (*Transaction, error) {
	return m.createFunc(t)
}
func (m *mockRepo) FindByExternalID(externalID string) (*Transaction, error) {
	return m.findByExternalIDFunc(externalID)
}
func (m *mockRepo) FindByID(id string) (*Transaction, error) {
	return m.findByIDFunc(id)
}
func (m *mockRepo) UpdateStatus(id, status string) error {
	return m.updateStatusFunc(id, status)
}
func (m *mockRepo) FindByUser(userID string) ([]Transaction, error) {
	return m.findByUserFunc(userID)
}

// --- mock PaymentGateway ---

type mockGateway struct {
	createInvoiceFunc func(externalID string, amount int64, description string) (*Invoice, error)
}

func (m *mockGateway) CreateInvoice(externalID string, amount int64, description string) (*Invoice, error) {
	return m.createInvoiceFunc(externalID, amount, description)
}

// --- mock PropertyChecker ---

type mockProperty struct {
	getOwnerAndStatusFunc func(propertyID string) (string, string, error)
	setStatusFunc         func(propertyID, status string) error
}

func (m *mockProperty) GetOwnerAndStatus(propertyID string) (string, string, error) {
	return m.getOwnerAndStatusFunc(propertyID)
}
func (m *mockProperty) SetStatus(propertyID, status string) error {
	return m.setStatusFunc(propertyID, status)
}