package status

import (
	"errors"
	"testing"
)

func TestUnfulfilled_String(t *testing.T) {
	u := unfulfilled{}
	if u.String() != string(Unfulfilled) {
		t.Errorf("Expected Unfulfilled, got %s", u.String())
	}
}

func TestUnfulfilled_Collected(t *testing.T) {
	u := unfulfilled{}
	_, err := u.Collected()
	if !errors.Is(err, AlreadyUnfulfilledError) {
		t.Errorf("Expected AlreadyUnfulfilledError, got %s", err)
	}
}

func TestUnfulfilled_Unfulfilled(t *testing.T) {
	u := unfulfilled{}
	_, err := u.Unfulfilled()
	if !errors.Is(err, AlreadyUnfulfilledError) {
		t.Errorf("Expected AlreadyUnfulfilledError, got %s", err)
	}
}

func TestUnfulfilled_Fulfilled(t *testing.T) {
	u := unfulfilled{}
	newState, err := u.Fulfilled()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if newState.Name() != Fulfilled {
		t.Errorf("Expected Fulfilled, got %s", u.Name())
	}
}
