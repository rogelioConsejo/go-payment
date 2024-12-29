package payment

import (
	"errors"
	"github.com/rogelioConsejo/go-payment/payment/status"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("It should return a payment with the method and status set", func(t *testing.T) {
		onCollect := func() error {
			return nil
		}
		p, _ := New("test", onCollect)
		if p.Method() != "test" {
			t.Errorf("Method() = %s; want test", p.Method())
		}
		if p.Status() != status.Pending {
			t.Errorf("Status() = %s; want %s", p.Status(), status.Pending)
		}
	})
	t.Run("It should return an error if onCollect is nil", func(t *testing.T) {
		p, err := New("test", nil)
		if p != nil {
			t.Errorf("New() = %v; want nil", p)
		}
		if !errors.Is(err, onCollectIsNilError) {
			t.Errorf("New() = %v; want %v", err, onCollectIsNilError)
		}
	})
}

func TestPayment_Fulfill(t *testing.T) {
	t.Run("It should call the onCollect function", func(t *testing.T) {
		called := false
		onCollect := func() error {
			called = true
			return nil
		}
		p, _ := New("test", onCollect)
		err := p.Fulfill()
		if err != nil {
			t.Errorf("Fulfill() = %v; want nil", err)
		}
		if !called {
			t.Error("onCollect was not called")
		}
	})
	t.Run("It should set the status to fulfilled", func(t *testing.T) {
		onCollect := func() error {
			return nil
		}
		p, _ := New("test", onCollect)
		_ = p.Fulfill()
		if p.Status() != status.Fulfilled {
			t.Errorf("Status() = %s; want %s", p.Status(), status.Fulfilled)
		}
	})
	t.Run("It should set the status to unfulfilled if onCollect fails", func(t *testing.T) {
		onCollect := func() error {
			return errors.New("error")
		}
		p, _ := New("test", onCollect)
		_ = p.Fulfill()
		if p.Status() != status.Unfulfilled {
			t.Errorf("Status() = %s; want %s", p.Status(), status.Unfulfilled)
		}
	})
	t.Run("It should return an error if the status is already fulfilled", func(t *testing.T) {
		onCollect := func() error {
			return nil
		}
		p, _ := New("test", onCollect)
		_ = p.Fulfill()
		err := p.Fulfill()
		if !errors.Is(err, status.AlreadyFulfilledError) {
			t.Errorf("Fulfill() = %v; want %v", err, status.AlreadyFulfilledError)
		}
	})
}
