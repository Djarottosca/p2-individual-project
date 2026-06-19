package transaction

import (
	"errors"

	"p2-individual-project/service/transaction"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(t *transaction.Transaction) (*transaction.Transaction, error) {
	if err := r.db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TransactionRepository) FindByExternalID(externalID string) (*transaction.Transaction, error) {
	var t transaction.Transaction
	err := r.db.Where("external_id = ?", externalID).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) FindByID(id string) (*transaction.Transaction, error) {
	var t transaction.Transaction
	err := r.db.Where("id = ?", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) UpdateStatus(id, status string) error {
	return r.db.Table("transactions").Where("id = ?", id).Update("status", status).Error
}

func (r *TransactionRepository) FindByUser(userID string) ([]transaction.Transaction, error) {
	var ts []transaction.Transaction
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&ts).Error
	if err != nil {
		return nil, err
	}
	return ts, nil
}
