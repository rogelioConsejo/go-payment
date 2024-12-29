package status

import "errors"

type pending struct {
}

func (p pending) String() string {
	return string(Pending)
}

func (p pending) Name() Name {
	return Pending
}

func (p pending) Collected() (State, error) {
	return collected{}, nil
}

func (p pending) Unfulfilled() (State, error) {
	return pending{}, NotCollectedError
}

func (p pending) Fulfilled() (State, error) {
	return pending{}, NotCollectedError
}

var NotCollectedError = errors.New("payment cannot be fulfilled if it has not been collected")
