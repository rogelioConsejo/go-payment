package payment

import "errors"

// Method should:
//   - Validate a Payment
//   - Initiate a Payment (and optionally, pre-authorize)
//   - Capture a Payment
//   - Refund a Payment
//   - Query a Payment status
type Method interface {
	Validate(Payment) error
	Create(Payment) (id string, err error)
}

var MethodIsNilError = errors.New("payment Method is nil")
