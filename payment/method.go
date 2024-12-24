package payment

import "errors"

type Method interface {
	Validate(Payment) error
	// Create may pre-authorize the payment if needed, but it should not charge the payment yet.
	Create(Payment) (ID, error)
	// Capture confirms the payment. Some payment methods may require some sort of validation to confirm that
	//the payment was in fact captured by the payment provider, like a secret nonce or hash.
	Capture(ID, Validation) error
}

type Validation string

var MethodIsNilError = errors.New("payment Method is nil")
