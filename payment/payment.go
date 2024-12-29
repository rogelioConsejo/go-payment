package payment

import (
	"errors"
	"github.com/rogelioConsejo/go-payment/payment/status"
)

type Payment interface {
	Method() MethodName
	Status() status.Name
	Fulfill() error
}

func New(method string, onCollect func() error) (Payment, error) {
	if onCollect == nil {
		return nil, onCollectIsNilError
	}
	return payment{
		method:           MethodName(method),
		status:           status.New(),
		executeAgreement: onCollect,
	}, nil
}

type ID string
type MethodName string

type payment struct {
	method           MethodName
	status           status.Status
	executeAgreement func() error
}

func (p payment) Fulfill() error {
	statusError := p.status.Collect()
	if statusError != nil {
		return errors.Join(CollectionStatusError, statusError)
	}

	if err := p.executeAgreement(); err != nil {
		statusChangeError := p.status.Unfulfill()
		if statusChangeError != nil {
			return errors.Join(ExecuteAgreementError, statusChangeError, err)
		}
		return errors.Join(ExecuteAgreementError, err)
	}
	if err := p.status.Fulfill(); err != nil {
		return errors.Join(FulfilledStatusError, err)
	}

	return nil
}

func (p payment) Status() status.Name {
	return status.Name(p.status.String())
}

func (p payment) Method() MethodName {
	return p.method
}

var IsNilError = errors.New("payment is nil")
var onCollectIsNilError = errors.New("onCollect callback cannot be nil when creating a new payment")
var ExecuteAgreementError = errors.New("could not execute agreement")
var CollectionStatusError = errors.New("could not change payment status to collected")
var FulfilledStatusError = errors.New("could not change payment status to fulfilled")
