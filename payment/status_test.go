package payment

import "testing"

func TestNewStatus(t *testing.T) {
	t.Run("NewStatus should return a status with the pending status", func(t *testing.T) {
		st := NewStatus()
		if st.String() != Pending {
			t.Errorf("expected status to be %s, got %s", Pending, st.String())
		}
	})
}

func TestStatus_Collect(t *testing.T) {
	t.Run("Collecting a payment should change its status to collected", func(t *testing.T) {
		st := NewStatus()
		st.Collect()
		if st.String() != Collected {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
	t.Run("Collecting a payment twice should not change its status", func(t *testing.T) {
		st := NewStatus()
		st.Collect()
		st.Collect()
		if st.String() != Collected {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
}
