# About
This service handle SCB internet banking payment by using Omise as a payment gateway.

## Setup
Create .env file by the running command below then edit Omise's credentials.
```
cp .env.example .env
```

## Run the app
Start the service on port 8080 (or set it via ```PORT``` variable).
```
make run
```

## Run unit tests
```
make test
```

## Build docker image
```
make docker-build
```

## Run docker container
```
docker run --env-file .env -p 8080:8080 paymentsvc
```

## Demo
Create a new payment request.
```
# Create a payment request of 2000 Satangs (20 THB)
curl -X POST http://localhost:8080/payments -d \
'{
    "amount": 2000,
    "currency": "THB",
    "return_uri": "https://example.com",
    "source_type": "internet_banking_scb"
}'
```
Response
```
{
    "id": 1,
    "authorized_uri": "https://pay.omise.co/offsites/ofsp_test_5mtr3e40dnsray0sxuk/pay"
}
```

Open a link in the ```authorized_uri``` field on the web browser then proceed to approve or reject the payment. The web browser will redirect to the ```return_uri``` specify on the first request.

Get the payment result.
```
# This will get a result of the payment with payment id = 1
curl http://localhost:8080/payments/1
```
Response
```
{
    "id": 1,
    "status": "successful",
    "amount": 2000,
    "currency": "THB",
    "source_type": "internet_banking_scb",
    "created_at": "2021-02-11T03:16:43.047466+07:00",
    "updated_at": "2021-02-11T03:33:29.65266+07:00"
}
```