package payment

import (
	"time"
)

// Service provides payment service methods.
type Service interface {
	CreatePaymentRequest(req *Request) (*Payment, error)
	Find(id int) (*Payment, error)
}

// Payment represents a payment.
type Payment struct {
	ID       int
	Status   Status
	Amount   int64
	Currency string

	OmiseCharge *OmiseCharge

	CreatedAt time.Time
	UpdatedAt time.Time
}

// OmiseCharge represents the Omise charge object.
// It should contains all attributes in Omise charge object https://www.omise.co/charges-api.
// But for the simplicity of this project we will use only the necessary attributes.
type OmiseCharge struct {
	ID           string
	Status       Status
	Amount       int64
	Currency     string
	AuthorizeURI string
	SourceType   string
	ReturnURI    string
}

// Status represents a payment status.
type Status string

// Payment statuses
const (
	StatusFailed     = "failed"
	StatusExpired    = "expired"
	StatusPending    = "pending"
	StatusReversed   = "reversed"
	StatusSuccessful = "successful"
)

// Repository provides access a data source.
type Repository interface {
	Create(payment *Payment) error
	Find(id int) (*Payment, error)
	UpdateStatus(id int, status Status) error
}

// Request contains details for making a payment.
type Request struct {
	Amount     int64
	Currency   string
	ReturnURI  string
	SourceType string
}

// Client provides methods for a payment gateway client to be implemented.
type Client interface {
	Charge(req *Request) (*OmiseCharge, error)
	GetCharge(id string) (*OmiseCharge, error)
}

type service struct {
	client Client
	repo   Repository
}

// NewService returns a new payment serivce.
func NewService(client Client, repo Repository) Service {
	return &service{
		client: client,
		repo:   repo,
	}
}

// CreatePaymentRequest creates a new payment request.
func (s *service) CreatePaymentRequest(req *Request) (*Payment, error) {
	charge, err := s.client.Charge(req)
	if err != nil {
		return nil, err
	}

	payment := &Payment{
		Status:      charge.Status,
		Amount:      charge.Amount,
		Currency:    charge.Currency,
		OmiseCharge: charge,
	}

	if err = s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// Find finds a payment with the given payment id in the data source.
// If payment status is pending, it will fetch for the updated payment through the payment client
// and store it in the data source.
func (s *service) Find(id int) (*Payment, error) {
	payment, err := s.repo.Find(id)
	if err != nil {
		return nil, err
	}

	if payment.OmiseCharge.Status != StatusPending {
		return payment, nil
	}

	charge, err := s.client.GetCharge(payment.OmiseCharge.ID)
	if err != nil {
		return nil, err
	}

	if err = s.repo.UpdateStatus(id, charge.Status); err != nil {
		return nil, err
	}

	payment, err = s.repo.Find(id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
