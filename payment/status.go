package payment

import "errors"

type Status interface {
	String() string
	Name() StatusName
	Collected() error
	Unfulfilled() error
	Fulfilled() error
}

type StatusName string

const (
	Pending     StatusName = "pending"
	Collected   StatusName = "collected"
	Unfulfilled StatusName = "unfulfilled"
	Fulfilled   StatusName = "fulfilled"
)

type status struct {
	current StatusName
}

func (s *status) Fulfilled() error {
	if s.current == Pending {
		return NotCollectedError
	}
	if s.current == Fulfilled {
		return AlreadyFulfilledError
	}
	s.current = Fulfilled
	return nil
}

func (s *status) Unfulfilled() error {
	if s.current == Pending {
		return NotCollectedError
	}
	if s.current == Fulfilled {
		return AlreadyFulfilledError
	}
	if s.current == Unfulfilled {
		return AlreadyUnfulfilledError
	}
	s.current = Unfulfilled
	return nil
}

func (s *status) Name() StatusName {
	return s.current
}

func (s *status) Collected() error {
	if s.current == Collected {
		return AlreadyCollectedError
	}
	s.current = Collected
	return nil
}

func (s *status) String() string {
	return string(s.current)
}

func NewStatus() Status {
	return &status{
		current: Pending,
	}
}

var AlreadyCollectedError = errors.New("payment has already been collected")
var NotCollectedError = errors.New("payment cannot be fulfilled if it has not been collected")
var AlreadyFulfilledError = errors.New("payment has already been fulfilled")
var AlreadyUnfulfilledError = errors.New("payment has already been unfulfilled")
