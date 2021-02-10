package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/noppawitt/paymentsvc/inmem"
	"github.com/noppawitt/paymentsvc/payment"
)

// Payment represents a payment handler.
type Payment struct {
	service payment.Service
}

// NewPayment returns a new payment handler.
func NewPayment(service payment.Service) *Payment {
	return &Payment{
		service: service,
	}
}

// Append appends routes to the router.
func (h *Payment) Append(r *mux.Router) {
	r.HandleFunc("", h.createPaymentRequest).Methods(http.MethodPost)
	r.HandleFunc("/{id}", h.getPayment).Methods(http.MethodGet)
}

type createPaymentRequestRequest struct {
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	ReturnURI  string `json:"return_uri"`
	SourceType string `json:"source_type"`
}

type createPaymentRequestResponse struct {
	ID            int    `json:"id"`
	AuthorizedURI string `json:"authorized_uri"`
}

func (h *Payment) createPaymentRequest(w http.ResponseWriter, r *http.Request) {
	req := &createPaymentRequestRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	paymentReq := &payment.Request{
		Amount:     req.Amount,
		Currency:   strings.ToUpper(req.Currency),
		ReturnURI:  req.ReturnURI,
		SourceType: req.SourceType,
	}

	payment, err := h.service.CreatePaymentRequest(paymentReq)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &createPaymentRequestResponse{
		ID:            payment.ID,
		AuthorizedURI: payment.OmiseCharge.AuthorizeURI,
	}

	respondJSON(w, res, http.StatusOK)
}

type getPaymentResponse struct {
	ID         int            `json:"id"`
	Status     payment.Status `json:"status"`
	Amount     int64          `json:"amount"`
	Currency   string         `json:"currency"`
	SourceType string         `json:"source_type"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (h *Payment) getPayment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, "payment id must be a number", http.StatusBadRequest)
		return
	}

	payment, err := h.service.Find(id)
	if err != nil {
		code := http.StatusInternalServerError
		if err == inmem.ErrPaymentNotFound {
			code = http.StatusBadRequest
		}
		respondError(w, err.Error(), code)
		return
	}

	res := &getPaymentResponse{
		ID:         payment.ID,
		Status:     payment.Status,
		Amount:     payment.Amount,
		Currency:   payment.Currency,
		SourceType: payment.OmiseCharge.SourceType,
		CreatedAt:  payment.CreatedAt,
		UpdatedAt:  payment.UpdatedAt,
	}

	respondJSON(w, res, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, err string, code int) {
	respondJSON(w, map[string]interface{}{"message": err}, code)
}
