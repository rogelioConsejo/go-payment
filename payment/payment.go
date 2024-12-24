package payment

import "errors"

type Payment interface {
	Method() MethodName
}

func New(method string) Payment {
	return payment{method: MethodName(method)}
}

type payment struct {
	method MethodName
}

func (t payment) Method() MethodName {
	return t.method
}

type MethodName string

var IsNilError = errors.New("payment is nil")
