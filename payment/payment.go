package payment

import "errors"

type Payment interface {
	Method() MethodName
	Status() Status
}

func New(method string) Payment {
	return payment{
		method: MethodName(method),
		status: NewStatus(),
	}
}

type ID string
type MethodName string

type payment struct {
	method MethodName
	status Status
}

func (t payment) Status() Status {
	return t.status
}

func (t payment) Method() MethodName {
	return t.method
}

var IsNilError = errors.New("payment is nil")
