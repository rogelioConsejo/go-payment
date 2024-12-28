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

func TestStatus_Collected(t *testing.T) {
	t.Run("Collecting a payment should change its status to collected", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		if st.String() != string(Collected) {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
	t.Run("Collecting a payment twice should not change its status", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		_ = st.Collected()
		if st.String() != string(Collected) {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
	t.Run("It should return an error if the status has already been collected", func(t *testing.T) {
		st := NewStatus()
		err := st.Collected()
		if err != nil {
			t.Errorf("there should be no error when collecting a payment for the first time, got %v", err)
		}
		err = st.Collected()
		if !errors.Is(err, AlreadyCollectedError) {
			t.Errorf("expected error to be %v, got %v", AlreadyCollectedError, err)
		}
	})
}

func TestStatus_Unfulfilled(t *testing.T) {
	t.Run("A payment can be unfulfilled after it is collected", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		err := st.Unfulfilled()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if st.Name() != Unfulfilled {
			t.Error("expected status to be unfulfilled")
		}
	})
	t.Run("It should return an error if the payment is pending", func(t *testing.T) {
		st := NewStatus()
		err := st.Unfulfilled()
		if !errors.Is(err, NotCollectedError) {
			t.Errorf("expected error to be %v, got %v", NotCollectedError, err)
		}
	})
	t.Run("It should return an error if the payment is already unfulfilled", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		err := st.Unfulfilled()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		err = st.Unfulfilled()
		if !errors.Is(err, AlreadyUnfulfilledError) {
			t.Errorf("expected error to be %v, got %v", AlreadyUnfulfilledError, err)
		}
	})
	t.Run("It should return an error if the payment is already fulfilled", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		_ = st.Fulfilled()
		err := st.Unfulfilled()
		if !errors.Is(err, AlreadyFulfilledError) {
			t.Errorf("expected error to be %v, got %v", AlreadyFulfilledError, err)
		}
	})
}

func TestStatus_Fulfilled(t *testing.T) {
	t.Run("A payment should be fulfilled when it is collected", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		err := st.Fulfilled()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if st.Name() != Fulfilled {
			t.Errorf("expected status to be %s", Fulfilled)
		}
	})
	t.Run("It should return an error if the payment is not pending", func(t *testing.T) {
		st := NewStatus()
		err := st.Fulfilled()
		if !errors.Is(err, NotCollectedError) {
			t.Errorf("expected error to be %v, got %v", NotCollectedError, err)
		}
	})
	t.Run("It should return an error if the payment is already fulfilled", func(t *testing.T) {
		st := NewStatus()
		_ = st.Collected()
		err := st.Fulfilled()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		err = st.Fulfilled()
		if !errors.Is(err, AlreadyFulfilledError) {
			t.Errorf("expected error to be %v, got %v", AlreadyFulfilledError, err)
		}
	})

}
