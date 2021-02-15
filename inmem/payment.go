package inmem

import (
	"errors"
	"sync"
	"time"

	"github.com/noppawitt/paymentsvc/payment"
)

// ErrPaymentNotFound occurs when a payment with given id is not exists.
var ErrPaymentNotFound = errors.New("payment not found")

// PaymentRepository provides access an in-memory data source.
type PaymentRepository struct {
	currentID int
	m         map[int]*payment.Payment
	mu        sync.RWMutex
}

// NewPaymentRepository returns a new payment repository.
func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		m: make(map[int]*payment.Payment),
	}
}

// Create creates a payment.
func (r *PaymentRepository) Create(payment *payment.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentID = r.currentID + 1
	payment.ID = r.currentID
	now := time.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now
	r.m[r.currentID] = payment
	return nil
}

// Find finds a payment with the given id.
func (r *PaymentRepository) Find(id int) (*payment.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	payment, ok := r.m[id]
	if !ok {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

// UpdateStatus updates a payment status of a payment with the given payment id.
func (r *PaymentRepository) UpdateStatus(id int, status payment.Status) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	payment, ok := r.m[id]
	if !ok {
		return ErrPaymentNotFound
	}
	payment.Status = status
	payment.OmiseCharge.Status = status
	payment.UpdatedAt = time.Now()
	return nil
}
