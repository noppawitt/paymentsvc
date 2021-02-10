package handler

import "github.com/noppawitt/paymentsvc/payment"

type mockService struct {
	CreatePaymentRequestFn func(req *payment.Request) (*payment.Payment, error)
	FindFn                 func(id int) (*payment.Payment, error)
}

func (m *mockService) CreatePaymentRequest(req *payment.Request) (*payment.Payment, error) {
	return m.CreatePaymentRequestFn(req)
}

func (m *mockService) Find(id int) (*payment.Payment, error) {
	return m.FindFn(id)
}
