package status

import "errors"

type fulfilled struct {
}

func (f fulfilled) String() string {
	return string(Fulfilled)
}

func (f fulfilled) Name() Name {
	return Fulfilled
}

func (f fulfilled) Collected() (State, error) {
	return fulfilled{}, AlreadyFulfilledError
}

func (f fulfilled) Unfulfilled() (State, error) {
	return fulfilled{}, AlreadyFulfilledError
}

func (f fulfilled) Fulfilled() (State, error) {
	return fulfilled{}, AlreadyFulfilledError
}

var AlreadyFulfilledError = errors.New("payment has already been fulfilled")
