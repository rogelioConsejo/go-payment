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

	id, initiatePaymentErr := perf.Initiate(New("cash"))
	if initiatePaymentErr != nil {
		t.Fatalf("could not initiate payment: %v", initiatePaymentErr)
	}

	checker := NewStatusChecker(spyPersistence)

	t.Run("The initial status of a payment should be pending", func(t *testing.T) {
		st, err := checker.CheckPaymentStatus(id)
		if err != nil {
			t.Fatalf("could not check payment status: %v", err)
		}

		if st.String() != Pending {
			t.Errorf("expected status to be %s, got %s", Pending, st.String())
		}
	})

	t.Run("The status of a payment should be collected after it is collected", func(t *testing.T) {
		err := perf.Confirm(id)
		if err != nil {
			t.Fatalf("could not confirm payment: %v", err)
		}
		st, err := checker.CheckPaymentStatus(id)
		if err != nil {
			t.Fatalf("could not check payment status: %v", err)
		}

		if st.String() != Collected {
			t.Errorf("expected status to be %s, got %s", Collected, st.String())
		}
	})
}

var _ StatusChecker = statusChecker{} // Ensure that statusChecker implements StatusChecker
