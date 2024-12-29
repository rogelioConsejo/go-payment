package status

import (
	"errors"
	"testing"
)

func TestFulfilled_String(t *testing.T) {
	f := fulfilled{}
	if f.String() != string(Fulfilled) {
		t.Errorf("expected status to be %s, got %s", Fulfilled, f.String())
	}
}

func TestFulfilled_Collected(t *testing.T) {
	t.Run("It should return an error if the status has already been fulfilled", func(t *testing.T) {
		f := fulfilled{}
		_, err := f.Collected()
		if !errors.Is(err, AlreadyFulfilledError) {
			t.Errorf("expected error to be %v, got %v", AlreadyFulfilledError, err)
		}
	})
}
