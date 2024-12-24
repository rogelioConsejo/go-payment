package payment

type Status interface {
	String() string
	Collect()
}

const (
	Pending   = "pending"
	Collected = "collected"
)

type status struct {
	current string
}

func (s *status) Collect() {
	s.current = Collected
}

func (s *status) String() string {
	return s.current
}

func NewStatus() Status {
	return &status{
		current: Pending,
	}
}
