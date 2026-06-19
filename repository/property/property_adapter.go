package property

import "p2-individual-project/service/property"

// Adapter ngebungkus property.Service biar memenuhi kontrak
// PropertyChecker yang dibutuhin service transaction.
type Adapter struct {
	svc property.Service
}

func NewAdapter(svc property.Service) *Adapter {
	return &Adapter{svc: svc}
}

func (a *Adapter) GetOwnerAndStatus(propertyID string) (string, string, error) {
	p, err := a.svc.Detail(propertyID)
	if err != nil {
		// properti nggak ada, balikin kosong tanpa bikin error fatal
		return "", "", nil
	}
	return p.UserID, p.Status, nil
}

func (a *Adapter) SetStatus(propertyID, status string) error {
	return a.svc.SetStatus(propertyID, status)
}
