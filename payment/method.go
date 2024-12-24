package payment

import "errors"

type Method interface {
	Validate(Payment) error
	// Create may pre-authorize the payment if needed, but it should not charge the payment yet.
	Create(Payment) (id string, err error)
}

var MethodIsNilError = errors.New("payment Method is nil")
