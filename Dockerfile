FROM golang:1.15-alpine3.13 AS builder
WORKDIR /go/src/github.com/noppawitt/paymentsvc
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o paymentsvc main.go

FROM scratch
WORKDIR /paymentsvc
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/noppawitt/paymentsvc ./
EXPOSE 8080
CMD [ "./paymentsvc" ]
