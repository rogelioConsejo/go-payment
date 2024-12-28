package payment

import (
	"errors"
	"testing"
)

func TestNewStatus(t *testing.T) {
	t.Run("NewStatus should return a status with the pending status", func(t *testing.T) {
		st := NewStatus()
		if st.String() != string(Pending) {
			t.Errorf("expected status to be %s, got %s", Pending, st.String())
		}
	})
}

func TestStatus_Collect(t *testing.T) {
	t.Run("Collecting a payment should change its status to collected", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collect()
		if st.String() != string(Collected) {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
	t.Run("Collecting a payment twice should not change its status", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collect()
		_ = st.Collect()
		if st.String() != string(Collected) {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
	t.Run("It should return an error if the status has already been collected", func(t *testing.T) {
		st := NewStatus()
		err := st.Collect()
		if err != nil {
			t.Errorf("there should be no error when collecting a payment for the first time, got %v", err)
		}
		err = st.Collect()
		if !errors.Is(err, AlreadyCollectedError) {
			t.Errorf("expected error to be %v, got %v", AlreadyCollectedError, err)
		}
	})
}
