package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/noppawitt/paymentsvc/client"
	"github.com/noppawitt/paymentsvc/handler"
	"github.com/noppawitt/paymentsvc/inmem"
	"github.com/noppawitt/paymentsvc/payment"
)

const defaultPort = "8080"

func main() {
	omisePublicKey := mustGetEnv("OMISE_PUBLIC_KEY")
	omiseSecretKey := mustGetEnv("OMISE_SECRET_KEY")
	port := getEnv("PORT", defaultPort)

	client := client.NewOmise(omisePublicKey, omiseSecretKey)

	paymentRepo := inmem.NewPaymentRepository()

	paymentSvc := payment.NewService(client, paymentRepo)

	paymentHandler := handler.NewPayment(paymentSvc)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Payment Service"))
	})

	paymentRouter := router.PathPrefix("/payments").Subrouter()
	paymentHandler.Append(paymentRouter)

	log.Println("Server is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func mustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(key + " is not defined")
	}
	return val
}
