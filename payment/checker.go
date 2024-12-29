package payment

import (
	"errors"
	"github.com/rogelioConsejo/go-payment/payment/status"
)

func NewStatusChecker(persistence RetrieverPersistence) StatusChecker {
	return statusChecker{
		persistence,
	}
}

type StatusChecker interface {
	CheckPaymentStatus(ID) (status.Name, error)
}

type statusChecker struct {
	RetrieverPersistence
}

func (s statusChecker) CheckPaymentStatus(id ID) (status.Name, error) {
	pay, err := s.RetrievePayment(string(id))
	if err != nil {
		return "", errors.Join(CheckPaymentStatusError, err)
	}
	return pay.Status(), nil
}

var CheckPaymentStatusError = errors.New("could not check payment status")
