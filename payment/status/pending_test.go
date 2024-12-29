package status

import "testing"

func TestPending_Name(t *testing.T) {
	p := pending{}
	if p.Name() != Pending {
		t.Errorf("Expected Pending, got %s", p.Name())
	}
}
