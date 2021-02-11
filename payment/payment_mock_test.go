package payment

type mockClient struct {
	ChargeFn    func(req *Request) (*OmiseCharge, error)
	GetChargeFn func(id string) (*OmiseCharge, error)
}

func (m *mockClient) Charge(req *Request) (*OmiseCharge, error) {
	return m.ChargeFn(req)
}

func (m *mockClient) GetCharge(id string) (*OmiseCharge, error) {
	return m.GetChargeFn(id)
}

type mockRepository struct {
	CreateFn        func(payment *Payment) error
	FindFn          func(id int) (*Payment, error)
	FindCalledTimes int
	UpdateStatusFn  func(id int, status Status) error
}

func (m *mockRepository) Create(payment *Payment) error {
	return m.CreateFn(payment)
}

func (m *mockRepository) Find(id int) (*Payment, error) {
	m.FindCalledTimes++
	return m.FindFn(id)
}

func (m *mockRepository) UpdateStatus(id int, status Status) error {
	return m.UpdateStatusFn(id, status)
}
