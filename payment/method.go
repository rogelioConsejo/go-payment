package payment

import (
	"errors"
)

// Method is the interface that all payment methods must implement.
type Method interface {
	// Validate checks if the payment is valid and can be processed by the payment method (amount, currency, etc).
	Validate(Payment) error
	// Create may pre-authorize the payment if needed, but it should not charge the payment yet.
	Create(Payment) (ID, error)
	// Capture confirms the payment. Some payment methods may require some sort of validation to confirm that
	//the payment was in fact captured by the payment provider, like a secret nonce or hash.
	Capture(ID, Validation) error
}

type Validation string

var MethodIsNilError = errors.New("payment Method is nil")
