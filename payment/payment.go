package payment

import "errors"

type Payment interface {
	Method() MethodName
	Status() StatusName
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

func (p payment) Fulfill() error {
	statusError := p.status.Collect()
	if statusError != nil {
		return errors.Join(CollectionStatusError, statusError)
	}
	if err := p.executeAgreement(); err != nil {
		return errors.Join(ExecuteAgreementError, err)
	}
	return nil
}

func (p payment) Status() StatusName {
	return StatusName(p.status.String())
}

func (p payment) Method() MethodName {
	return p.method
}

var IsNilError = errors.New("payment is nil")
var onCollectIsNilError = errors.New("onCollect callback cannot be nil when creating a new payment")
var ExecuteAgreementError = errors.New("could not execute agreement")
var CollectionStatusError = errors.New("could not change payment status to collected")
