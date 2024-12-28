package payment

import "errors"

type Status interface {
	String() string
	Name() StatusName
	Collect() error
}

type StatusName string

const (
	Pending   StatusName = "pending"
	Collected StatusName = "collected"
)

type status struct {
	current StatusName
}

func (s *status) Name() StatusName {
	return s.current
}

func (s *status) Collect() error {
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
