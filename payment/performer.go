package payment

import "errors"

type Performer interface {
	PerformerInitializer
	Initiate(Payment) (ID, error)
	Confirm(ID, Validation) error
}

type PerformerInitializer interface {
	AddPaymentMethod(MethodName, Method) error
}

func NewPaymentPerformer(p PerformerPersistence) Performer {
	return performer{p}
}

type PerformerPersistence interface {
	SavePaymentMethod(MethodName, Method) error
	RetrievePaymentMethod(MethodName) (Method, error)
	SavePayment(id string, pay Payment) error
	RetrieverPersistence
}

type RetrieverPersistence interface {
	RetrievePayment(id string) (Payment, error)
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

func (p performer) Initiate(payment Payment) (id ID, err error) {
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
	saveErr := p.SavePayment(string(id), payment)
	if saveErr != nil {
		return "", errors.Join(SaveError, saveErr)
	}

	return id, nil
}

func (p performer) Confirm(id ID, v Validation) error {
	if id == "" {
		return EmptyPaymentIDError
	}
	if p.PerformerPersistence == nil {
		return PersistenceNotSetError
	}
	pay, err := p.RetrievePayment(string(id))
	if err != nil {
		return NotFoundError
	}
	method, err := p.RetrievePaymentMethod(pay.Method())
	if err != nil {
		return errors.Join(UnsupportedMethodError, MethodMayHaveBeenRemovedError) // This should never happen because the payment was already validated on creation, but just in case
	}

	err = method.Capture(id, v)
	if err != nil {
		return errors.Join(CaptureError, err)
	}

	pay.Status().Collect()
	fulfillmentErr := pay.Fulfill()
	if fulfillmentErr != nil {
		return errors.Join(FulfillmentError, fulfillmentErr)
	}
	err = p.SavePayment(string(id), pay)
	if err != nil {
		return errors.Join(SaveError, err)
	}

	return nil
}

var EmptyMethodError = errors.New("payment MethodName is empty")
var UnsupportedMethodError = errors.New("payment MethodName is not supported")
var PersistenceNotSetError = errors.New("persistence is not set")
var InvalidPaymentError = errors.New("payment is invalid for the selected payment method")
var CreationError = errors.New("payment creation failed")
var SaveError = errors.New("payment save failed")
var SaveMethodError = errors.New("payment method save failed")
var EmptyPaymentIDError = errors.New("payment ID is empty")
var NotFoundError = errors.New("payment not found")
var CaptureError = errors.New("payment capture failed")
var FulfillmentError = errors.New("payment fulfillment failed")

// MethodMayHaveBeenRemovedError exists in case the payment method was removed after the payment was created and before
// it was confirmed, which would be a very bad operative practice, but it could lead to a lot of confusion, so
// it's better to be explicit about it, even if it is unlikely to happen.
var MethodMayHaveBeenRemovedError = errors.New("payment method may have been removed before payment was confirmed")
