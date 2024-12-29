package status

import "errors"

type collected struct {
}

func (c collected) String() string {
	return string(Collected)
}

func (c collected) Name() Name {
	return Collected
}

func (c collected) Collected() (State, error) {
	return collected{}, AlreadyCollectedError
}

func (c collected) Unfulfilled() (State, error) {
	return unfulfilled{}, nil
}

func (c collected) Fulfilled() (State, error) {
	return fulfilled{}, nil
}

var AlreadyCollectedError = errors.New("payment has already been collected")
