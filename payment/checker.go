package payment

import "errors"

func NewStatusChecker(persistence RetrieverPersistence) StatusChecker {
	return statusChecker{
		persistence,
	}
}

type StatusChecker interface {
	CheckPaymentStatus(ID) (Status, error)
}

type statusChecker struct {
	RetrieverPersistence
}

func (s statusChecker) CheckPaymentStatus(id ID) (Status, error) {
	pay, err := s.RetrievePayment(string(id))
	if err != nil {
		return nil, errors.Join(CheckPaymentStatusError, err)
	}
	return pay.Status(), nil
}

var CheckPaymentStatusError = errors.New("could not check payment status")
