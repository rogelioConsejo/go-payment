package payment

import "errors"

type Performer interface {
	AddPaymentMethod(MethodName, Method) error
	Initiate(Payment) (id string, err error)
}

func NewPaymentPerformer(p PerformerPersistence) Performer {
	return performer{p}
}

type PerformerPersistence interface {
	SavePaymentMethod(MethodName, Method) error
	RetrievePaymentMethod(MethodName) (Method, error)
	SavePayment(id string, pay Payment) error
}

type performer struct {
	PerformerPersistence
}

func (p performer) AddPaymentMethod(name MethodName, method Method) error {
	if p.PerformerPersistence == nil {
		return PersistenceNotSetError
	}
	if name == "" {
		return EmptyMethodError
	}
	if method == nil {
		return MethodIsNilError
	}
	err := p.SavePaymentMethod(name, method)
	if err != nil {
		return errors.Join(SaveMethodError, err)
	}
	return nil
}

func (p performer) Initiate(payment Payment) (id string, err error) {
	if payment == nil {
		return "", IsNilError
	}
	if payment.Method() == "" {
		return "", EmptyMethodError
	}
	method, err := p.RetrievePaymentMethod(payment.Method())
	if err != nil {
		return "", UnsupportedMethodError
	}
	validationErr := method.Validate(payment)
	if validationErr != nil {
		return "", errors.Join(InvalidPaymentError, validationErr)
	}
	id, creationError := method.Create(payment)
	if creationError != nil {
		return "", errors.Join(CreationError, creationError)
	}
	saveErr := p.SavePayment(id, payment)
	if saveErr != nil {
		return "", errors.Join(SaveError, saveErr)
	}

	return id, nil
}

var EmptyMethodError = errors.New("payment MethodName is empty")
var UnsupportedMethodError = errors.New("payment MethodName is not supported")
var PersistenceNotSetError = errors.New("persistence is not set")
var InvalidPaymentError = errors.New("payment is invalid for the selected payment method")
var CreationError = errors.New("payment creation failed")
var SaveError = errors.New("payment save failed")
var SaveMethodError = errors.New("payment method save failed")
