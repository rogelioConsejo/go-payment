package payment

import "errors"

type Payment interface {
	Method() MethodName
	Status() Status
	Fulfill() error
}

func New(method string, onCollect func() error) (Payment, error) {
	if onCollect == nil {
		return nil, onCollectIsNilError
	}
	return payment{
		method:           MethodName(method),
		status:           NewStatus(),
		executeAgreement: onCollect,
	}, nil
}

type ID string
type MethodName string

type payment struct {
	method           MethodName
	status           Status
	executeAgreement func() error
}

func (t payment) Fulfill() error {
	return t.executeAgreement()
}

func (t payment) Status() Status {
	return t.status
}

func (t payment) Method() MethodName {
	return t.method
}

var IsNilError = errors.New("payment is nil")
var onCollectIsNilError = errors.New("onCollect is nil")
