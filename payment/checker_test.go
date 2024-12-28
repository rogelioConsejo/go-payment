package payment

import "testing"

func TestNewStatusChecker(t *testing.T) {
	var checker StatusChecker = NewStatusChecker(getSpyPerformerPersistence())
	if checker == nil {
		t.Error("NewStatusChecker should not return nil")
	}
}

func TestStatusChecker_CheckPaymentStatus(t *testing.T) {
	spyPersistence := getSpyPerformerPersistence()
	perf := NewPaymentPerformer(spyPersistence)

	cashMethod := getSpyMethod()
	addPaymentErr := perf.AddPaymentMethod("cash", cashMethod)
	if addPaymentErr != nil {
		t.Fatalf("could not add payment method: %v", addPaymentErr)
	}

	pay, _ := New("cash", onCollectStub)
	id, initiatePaymentErr := perf.Initiate(pay)
	if initiatePaymentErr != nil {
		t.Fatalf("could not initiate payment: %v", initiatePaymentErr)
	}

	checker := NewStatusChecker(spyPersistence)

	t.Run("The initial status of a payment should be pending", func(t *testing.T) {
		st, err := checker.CheckPaymentStatus(id)
		if err != nil {
			t.Fatalf("could not check payment status: %v", err)
		}

		if st != Pending {
			t.Errorf("expected status to be %s, got %s", Pending, st)
		}
	})

	t.Run("The status of a payment should be fulfilled after it is confirmed", func(t *testing.T) {
		err := perf.Confirm(id, "valid")
		if err != nil {
			t.Fatalf("could not confirm payment: %v", err)
		}
		st, err := checker.CheckPaymentStatus(id)
		if err != nil {
			t.Fatalf("could not check payment status: %v", err)
		}

		if st != Fulfilled {
			t.Errorf("expected status to be %s, got %s", Fulfilled, st)
		}
	})
}

var _ StatusChecker = statusChecker{} // Ensure that statusChecker implements StatusChecker
