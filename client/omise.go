package client

import (
	"log"

	"github.com/noppawitt/paymentsvc/payment"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// Omise is a wrapper of the Omise Go client.
type Omise struct {
	client *omise.Client
}

// NewOmise returns a new omise client.
func NewOmise(publicKey, secretKey string) *Omise {
	client, err := omise.NewClient(publicKey, secretKey)
	if err != nil {
		log.Fatal(err)
	}
	return &Omise{
		client: client,
	}
}

// Charge charges the payment source.
func (c *Omise) Charge(req *payment.Request) (*payment.OmiseCharge, error) {
	source := &omise.Source{}
	createSource := &operations.CreateSource{
		Type:     req.SourceType,
		Amount:   req.Amount,
		Currency: req.Currency,
	}

	if err := c.client.Do(source, createSource); err != nil {
		return nil, err
	}

	charge := &omise.Charge{}
	createCharge := &operations.CreateCharge{
		Source:    source.ID,
		Amount:    source.Amount,
		Currency:  source.Currency,
		ReturnURI: req.ReturnURI,
	}

	if err := c.client.Do(charge, createCharge); err != nil {
		return nil, err
	}

	omiseCharge := &payment.OmiseCharge{
		ID:           charge.ID,
		Status:       payment.Status(charge.Status),
		Amount:       charge.Amount,
		Currency:     charge.Currency,
		AuthorizeURI: charge.AuthorizeURI,
		SourceType:   charge.Source.Type,
		ReturnURI:    charge.ReturnURI,
	}

	return omiseCharge, nil
}

// GetCharge gets a charge with the given charge id.
func (c *Omise) GetCharge(id string) (*payment.OmiseCharge, error) {
	charge := &omise.Charge{}
	retrieve := &operations.RetrieveCharge{ChargeID: id}
	if err := c.client.Do(charge, retrieve); err != nil {
		return nil, err
	}

	omiseCharge := &payment.OmiseCharge{
		ID:           charge.ID,
		Status:       payment.Status(charge.Status),
		Amount:       charge.Amount,
		Currency:     charge.Currency,
		AuthorizeURI: charge.AuthorizeURI,
		SourceType:   charge.Source.Type,
		ReturnURI:    charge.ReturnURI,
	}

	return omiseCharge, nil
}
