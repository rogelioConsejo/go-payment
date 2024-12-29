package status

import "errors"

type unfulfilled struct {
}

func (u unfulfilled) String() string {
	return string(Unfulfilled)
}

func (u unfulfilled) Name() Name {
	return Unfulfilled
}

func (u unfulfilled) Collected() (State, error) {
	return unfulfilled{}, AlreadyUnfulfilledError
}

func (u unfulfilled) Unfulfilled() (State, error) {
	return unfulfilled{}, AlreadyUnfulfilledError
}

func (u unfulfilled) Fulfilled() (State, error) {
	return fulfilled{}, nil
}

var AlreadyUnfulfilledError = errors.New("payment has already been collected but could not be fulfilled, it cannot be collected or marked as unfulfilled again")
