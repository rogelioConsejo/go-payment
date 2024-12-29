package status

import "errors"

func New() Status {
	return &status{
		pending{},
	}
}

type Status interface {
	StateQueries
	Collect() error
	Unfulfill() error
	Fulfill() error
}

type StateQueries interface {
	String() string
	Name() Name
}

type StateChanges interface {
	Collected() (State, error)
	Unfulfilled() (State, error)
	Fulfilled() (State, error)
}

type State interface {
	StateQueries
	StateChanges
}

type Name string

const (
	Pending     Name = "pending"
	Collected   Name = "collected"
	Unfulfilled Name = "unfulfilled"
	Fulfilled   Name = "fulfilled"
)

// maybe we can have structs for each status to make this easier to extend. They would need to implement StateChanges
// more like a proper state machine. Instead of having the "statusName", we could embed the state. This may also need
// to be a separate package for the state machine to be more flexible and clean.
type status struct {
	State
}

func (s *status) Collect() error {
	newState, err := s.Collected()
	if err != nil {
		return errors.Join(CollectionStatusError, err)
	}
	s.State = newState
	return nil
}

func (s *status) Unfulfill() error {
	newState, err := s.Unfulfilled()
	if err != nil {
		return errors.Join(UnfulfillmentStatusError, err)
	}
	s.State = newState
	return nil
}

func (s *status) Fulfill() error {
	newState, err := s.Fulfilled()
	if err != nil {
		return errors.Join(FulfilledStatusError, err)
	}
	s.State = newState
	return nil
}

var CollectionStatusError = errors.New("could not change payment status to collected")
var UnfulfillmentStatusError = errors.New("could not change payment status to unfulfilled")
var FulfilledStatusError = errors.New("could not change payment status to fulfilled")
