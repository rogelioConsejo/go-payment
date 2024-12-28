package payment

import "errors"

type Status interface {
	StateQueries
	StateChanges
}

type StateQueries interface {
	String() string
	Name() StatusName
}

type StateChanges interface {
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

// maybe we can have structs for each status to make this easier to extend. They would need to implement StateChanges
// more like a proper state machine. Instead of having the "statusName", we could embed the state. This may also need
// to be a separate package for the state machine to be more flexible and clean.
type status struct {
	current StatusName
}

func (s *status) Fulfilled() error {
	switch s.current {
	case Pending:
		return NotCollectedError
	case Fulfilled:
		return AlreadyFulfilledError
	case Unfulfilled:
		s.current = Fulfilled
	case Collected:
		s.current = Fulfilled
	default:
		return UnknownStatusError
	}

	return nil
}

func (s *status) Unfulfilled() error {
	switch s.current {
	case Pending:
		return NotCollectedError
	case Unfulfilled:
		return AlreadyUnfulfilledError
	case Fulfilled:
		return AlreadyFulfilledError
	case Collected:
		s.current = Unfulfilled
	default:
		return UnknownStatusError
	}

	return nil
}

func (s *status) Name() StatusName {
	return s.current
}

func (s *status) Collected() error {
	switch s.current {
	case Collected:
		return AlreadyCollectedError
	case Fulfilled:
		return AlreadyFulfilledError
	case Unfulfilled:
		return AlreadyUnfulfilledError
	case Pending:
		s.current = Collected
	default:
		return UnknownStatusError
	}

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
var AlreadyUnfulfilledError = errors.New("payment has already been collected but could not be fulfilled")
var UnknownStatusError = errors.New("unknown status")
