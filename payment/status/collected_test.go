package status

import "testing"

func TestCollected_Name(t *testing.T) {
	st := collected{}
	if st.Name() != Collected {
		t.Errorf("expected status to be %s, got %s", Collected, st.Name())
	}
}
