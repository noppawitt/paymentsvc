package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/noppawitt/paymentsvc/inmem"
	"github.com/noppawitt/paymentsvc/payment"
)

var (
	now        = time.Now()
	nowJSON, _ = now.MarshalJSON()
)

func paymentRouter() *mux.Router {
	return mux.NewRouter().PathPrefix("/payments").Subrouter()
}

func TestPayment_createPaymentRequest(t *testing.T) {
	tests := []struct {
		name                       string
		reqBody                    string
		createPaymentRequestReturn *payment.Payment
		createPaymentRequestErr    error
		want                       string
		wantStatus                 int
	}{
		{
			name:    "success",
			reqBody: `{"amount":2000,"currency":"THB","return_uri":"http://localhost:8080","payment_type":"internet_banking_scb"}`,
			createPaymentRequestReturn: &payment.Payment{
				ID:       1,
				Status:   payment.StatusPending,
				Amount:   2000,
				Currency: "THB",
				OmiseCharge: &payment.OmiseCharge{
					ID:           "charge-1",
					Status:       payment.StatusPending,
					Amount:       2000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			createPaymentRequestErr: nil,
			want:                    `{"id":1,"authorized_uri":"http://authuri.com"}`,
			wantStatus:              http.StatusOK,
		},
		{
			name:                       "invalid request body",
			reqBody:                    `x`,
			createPaymentRequestReturn: nil,
			createPaymentRequestErr:    nil,
			want:                       `{"message":"invalid request body"}`,
			wantStatus:                 http.StatusBadRequest,
		},
		{
			name:                       "error",
			reqBody:                    `{"amount":2000,"currency":"THB","return_uri":"http://localhost:8080","payment_type":"internet_banking_scb"}`,
			createPaymentRequestReturn: nil,
			createPaymentRequestErr:    errors.New("some error"),
			want:                       `{"message":"some error"}`,
			wantStatus:                 http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockService{}
			s.CreatePaymentRequestFn = func(req *payment.Request) (*payment.Payment, error) {
				return tt.createPaymentRequestReturn, tt.createPaymentRequestErr
			}

			r := paymentRouter()
			h := NewPayment(s)
			h.Append(r)

			req, err := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer([]byte(tt.reqBody)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if gotStatus := rr.Code; gotStatus != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", gotStatus, tt.wantStatus)
			}

			got := strings.TrimSpace(rr.Body.String())
			if got != tt.want {
				t.Errorf("handler returned unexpected body: got %v want %v", got, tt.want)
			}
		})
	}
}

func TestPayment_getPayment(t *testing.T) {
	tests := []struct {
		name       string
		paymentID  string
		FindReturn *payment.Payment
		FindErr    error
		want       string
		wantStatus int
	}{
		{
			name:      "success",
			paymentID: "1",
			FindReturn: &payment.Payment{
				ID:       1,
				Status:   payment.StatusSuccessful,
				Amount:   20000,
				Currency: "THB",
				OmiseCharge: &payment.OmiseCharge{
					ID:           "charge-1",
					Status:       payment.StatusSuccessful,
					Amount:       20000,
					Currency:     "THB",
					AuthorizeURI: "http://authuri.com",
					SourceType:   "internet_banking_scb",
					ReturnURI:    "http://returnuri.com",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			FindErr:    nil,
			want:       fmt.Sprintf(`{"id":1,"status":"successful","amount":20000,"currency":"THB","source_type":"internet_banking_scb","created_at":%s,"updated_at":%s}`, nowJSON, nowJSON),
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid payment id",
			paymentID:  "x",
			FindReturn: nil,
			FindErr:    inmem.ErrPaymentNotFound,
			want:       `{"message":"payment id must be a number"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			paymentID:  "1",
			FindReturn: nil,
			FindErr:    inmem.ErrPaymentNotFound,
			want:       `{"message":"payment not found"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error",
			paymentID:  "1",
			FindReturn: nil,
			FindErr:    errors.New("some error"),
			want:       `{"message":"some error"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockService{}
			s.FindFn = func(id int) (*payment.Payment, error) {
				return tt.FindReturn, tt.FindErr
			}

			r := paymentRouter()
			h := NewPayment(s)
			h.Append(r)

			req, err := http.NewRequest(http.MethodGet, "/payments/"+tt.paymentID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if gotStatus := rr.Code; gotStatus != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", gotStatus, tt.wantStatus)
			}

			got := strings.TrimSpace(rr.Body.String())
			if got != tt.want {
				t.Errorf("handler returned unexpected body: got %v want %v", got, tt.want)
			}
		})
	}
}
