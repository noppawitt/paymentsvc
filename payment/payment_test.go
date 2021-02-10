package payment

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

var (
	now          = time.Now()
	errSomeError = errors.New("some error")
)

func TestService_CreatePaymentRequest(t *testing.T) {
	type mocks struct {
		paymentStatus   Status
		clientReturnErr error
		paymentID       int
		repoReturnErr   error
		omiseChargeID   string
		authorizeURI    string
	}
	type args struct {
		req *Request
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *Payment
		wantErr bool
	}{
		{
			name: "success",
			mocks: mocks{
				paymentStatus:   StatusPending,
				clientReturnErr: nil,
				paymentID:       1,
				repoReturnErr:   nil,
				omiseChargeID:   "charge-1",
				authorizeURI:    "http://authuri.com",
			},
			args: args{
				req: &Request{
					Amount:     20000,
					Currency:   "THB",
					ReturnURI:  "http://returnuri.com",
					SourceType: "internet_banking_scb",
				},
			},
			want: &Payment{
				ID:       1,
				Status:   "pending",
				Amount:   20000,
				Currency: "THB",
				OmiseCharge: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusPending,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
		},
		{
			name: "client Error",
			mocks: mocks{
				paymentStatus:   StatusPending,
				clientReturnErr: errSomeError,
				paymentID:       0,
				repoReturnErr:   nil,
				omiseChargeID:   "",
				authorizeURI:    "",
			},
			args: args{
				req: &Request{
					Amount:     20000,
					Currency:   "THB",
					ReturnURI:  "http://returnuri.com",
					SourceType: "internet_banking_scb",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repository Error",
			mocks: mocks{
				paymentStatus:   StatusPending,
				clientReturnErr: nil,
				paymentID:       0,
				repoReturnErr:   errSomeError,
				omiseChargeID:   "charge-1",
				authorizeURI:    "http://authuri.com",
			},
			args: args{
				req: &Request{
					Amount:     20000,
					Currency:   "THB",
					ReturnURI:  "http://returnuri.com",
					SourceType: "internet_banking_scb",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockClient{}
			repo := &mockRepository{}

			client.ChargeFn = func(req *Request) (*OmiseCharge, error) {
				charge := &OmiseCharge{
					ID:           tt.mocks.omiseChargeID,
					Status:       tt.mocks.paymentStatus,
					Amount:       tt.args.req.Amount,
					Currency:     tt.args.req.Currency,
					AuthorizeURI: tt.mocks.authorizeURI,
					SourceType:   tt.args.req.SourceType,
					ReturnURI:    tt.args.req.ReturnURI,
				}
				return charge, tt.mocks.clientReturnErr
			}

			repo.CreateFn = func(payment *Payment) error {
				payment.ID = tt.mocks.paymentID
				payment.CreatedAt = now
				payment.UpdatedAt = now
				return tt.mocks.repoReturnErr
			}

			s := NewService(client, repo)
			got, err := s.CreatePaymentRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreatePaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreatePaymentRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Find(t *testing.T) {
	type mocks struct {
		findReturns     [2]*Payment
		findErrs        [2]error
		getChargeReturn *OmiseCharge
		getChargeErr    error
		updateStatusErr error
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *Payment
		wantErr error
	}{
		{
			name: "success not update status",
			mocks: mocks{
				findReturns: [2]*Payment{
					{
						ID:       1,
						Status:   StatusSuccessful,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusSuccessful,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				findErrs: [2]error{
					nil,
				},
				getChargeReturn: nil,
				getChargeErr:    nil,
				updateStatusErr: nil,
			},
			args: args{
				id: 1,
			},
			want: &Payment{
				ID:       1,
				Status:   StatusSuccessful,
				Amount:   20000,
				Currency: "THB",
				OmiseCharge: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "success update status",
			mocks: mocks{
				findReturns: [2]*Payment{
					{
						ID:       1,
						Status:   StatusPending,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusPending,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
					{
						ID:       1,
						Status:   StatusSuccessful,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusSuccessful,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				findErrs: [2]error{
					nil,
					nil,
				},
				getChargeReturn: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				getChargeErr:    nil,
				updateStatusErr: nil,
			},
			args: args{
				id: 1,
			},
			want: &Payment{
				ID:       1,
				Status:   StatusSuccessful,
				Amount:   20000,
				Currency: "THB",
				OmiseCharge: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "Find 1 error",
			mocks: mocks{
				findReturns: [2]*Payment{
					{
						ID:       1,
						Status:   StatusPending,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusPending,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
					{
						ID:       1,
						Status:   StatusSuccessful,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusSuccessful,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				findErrs: [2]error{
					errSomeError,
					nil,
				},
				getChargeReturn: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				getChargeErr:    nil,
				updateStatusErr: nil,
			},
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: errSomeError,
		},
		{
			name: "Find 2 error",
			mocks: mocks{
				findReturns: [2]*Payment{
					{
						ID:       1,
						Status:   StatusPending,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusPending,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
					{
						ID:       1,
						Status:   StatusSuccessful,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusSuccessful,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				findErrs: [2]error{
					nil,
					errSomeError,
				},
				getChargeReturn: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				getChargeErr:    nil,
				updateStatusErr: nil,
			},
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: errSomeError,
		},
		{
			name: "GetCharge error",
			mocks: mocks{
				findReturns: [2]*Payment{
					{
						ID:       1,
						Status:   StatusPending,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusPending,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				findErrs: [2]error{
					nil,
				},
				getChargeReturn: nil,
				getChargeErr:    errSomeError,
				updateStatusErr: nil,
			},
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: errSomeError,
		},
		{
			name: "UpdateStatus error",
			mocks: mocks{
				findReturns: [2]*Payment{
					{
						ID:       1,
						Status:   StatusPending,
						Amount:   20000,
						Currency: "THB",
						OmiseCharge: &OmiseCharge{
							ID:           "charge-1",
							Status:       StatusPending,
							Amount:       20000,
							Currency:     "THB",
							AuthorizeURI: "http://authuri.com",
							SourceType:   "internet_banking_scb",
							ReturnURI:    "http://returnuri.com",
						},
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				findErrs: [2]error{
					nil,
				},
				getChargeReturn: &OmiseCharge{
					ID:           "charge-1",
					Status:       StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				getChargeErr:    nil,
				updateStatusErr: errSomeError,
			},
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: errSomeError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockClient{}
			repo := &mockRepository{}

			client.GetChargeFn = func(id string) (*OmiseCharge, error) {
				return tt.mocks.getChargeReturn, tt.mocks.getChargeErr
			}

			repo.FindFn = func(id int) (*Payment, error) {
				return tt.mocks.findReturns[repo.FindCalledTimes-1], tt.mocks.findErrs[repo.FindCalledTimes-1]
			}

			repo.UpdateStatusFn = func(id int, status Status) error {
				return tt.mocks.updateStatusErr
			}

			s := NewService(client, repo)
			got, err := s.Find(tt.args.id)
			if err != tt.wantErr {
				t.Errorf("Service.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
