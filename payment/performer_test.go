package payment

import (
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestPerformer_AddPaymentMethod(t *testing.T) {
	per := getSpyPerformerPersistence()
	var paymentPerformer Performer = NewPaymentPerformer(per)
	if paymentPerformer == nil {
		t.Errorf("Expected Performer, got nil")
	}

	t.Run("Add payment method should save the payment method through persistence", func(t *testing.T) {
		_ = paymentPerformer.AddPaymentMethod("paypal", getSpyMethod())
		ok := per.paymentMethodExists("paypal")
		if !ok {
			t.Errorf("Expected payment method to be saved, got not saved")
		}
	})
	t.Run("It should return an error if the method name is empty", func(t *testing.T) {
		err := paymentPerformer.AddPaymentMethod("", nil)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, EmptyMethodError) {
			t.Errorf("Expected error to be %v, got %v", EmptyMethodError, err)
		}
		ok := per.paymentMethodExists("")
		if ok {
			t.Errorf("Expected payment method to not be saved, got saved")
		}
	})
	t.Run("It should return an error if the Method is nil", func(t *testing.T) {
		const methodName MethodName = "nil-method"
		err := paymentPerformer.AddPaymentMethod(methodName, nil)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, MethodIsNilError) {
			t.Errorf("Expected error to be %v, got %v", IsNilError, err)
		}
		ok := per.paymentMethodExists(methodName)
		if ok {
			t.Errorf("Expected payment method to not be saved, got saved")
		}
	})
	t.Run("It should return an error if the persistence fails to save the payment method", func(t *testing.T) {
		per.failWhenSavingAPaymentMethod(true)
		err := paymentPerformer.AddPaymentMethod("paypal", getSpyMethod())
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, SaveMethodError) {
			t.Errorf("Expected error to be %v, got %v", SaveMethodError, err)
		}
		per.failWhenSavingAPaymentMethod(false)
	})
}

func TestPerformer_Initiate(t *testing.T) {
	spyPersistence := getSpyPerformerPersistence()
	var per Performer = NewPaymentPerformer(spyPersistence)
	if per == nil {
		t.Fatalf("Expected Performer, got nil")
	}

	m := getSpyMethod()
	_ = per.AddPaymentMethod("paypal", m)

	t.Run("Initiate payment should return an error if the payment is nil", func(t *testing.T) {
		testNilPayment(t, per)
		testNonNilPayment(t, per)
	})
	t.Run("Initiate payment should return an error if the payment has no MethodName", func(t *testing.T) {
		testPaymentWithoutMethod(t, per)
		testPaymentWithMethod(t, per)
	})
	t.Run("Initiate payment should an error if the payment method is not supported", func(t *testing.T) {
		testUnsupportedMethod(t, per)
		testSupportedMethod(t, per)
	})
	t.Run("It should return a payment ID if the payment is valid", func(t *testing.T) {
		pay := New("paypal")
		id, err := per.Initiate(pay)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if id == "" {
			t.Errorf("Expected payment ID, got empty string")
		}
	})
	t.Run("It should validate the payment using the payment method", func(t *testing.T) {
		pay := New("paypal")
		_, _ = per.Initiate(pay)
		if !m.wasValidated(pay) {
			t.Errorf("Expected payment to be validated, got not validated")
		}
	})
	t.Run("It should return an error if the payment is invalid for the method", func(t *testing.T) {
		m.rejectAllPayments(true)
		pay := New("paypal")
		_, err := per.Initiate(pay)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, InvalidPaymentError) {
			t.Errorf("Expected error to be %v, got %v", InvalidPaymentError, err)
		}
		m.rejectAllPayments(false)
	})
	t.Run("It should initiate the payment via the payment method", func(t *testing.T) {
		pay := New("paypal")
		id, err := per.Initiate(pay)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if id == "" {
			t.Errorf("Expected payment ID, got empty string")
		}
		if !m.wasCreated(id) {
			t.Errorf("Expected payment to be initiated, got not initiated")
		}
	})
	t.Run("It should return an error if the payment could not be initiated", func(t *testing.T) {
		m.failAllPayments(true)
		pay := New("paypal")
		_, err := per.Initiate(pay)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, CreationError) {
			t.Errorf("Expected error to be %v, got %v", CreationError, err)
		}
		m.failAllPayments(false)
	})
	t.Run("It should save the created payment by its ID", func(t *testing.T) {
		pay := New("paypal")
		id, _ := per.Initiate(pay)
		if !spyPersistence.paymentExists(id) {
			t.Errorf("Expected payment to be saved, got not saved")
		}
	})
	t.Run("It should return an error if the payment could not be saved", func(t *testing.T) {
		spyPersistence.failWhenSavingAPayment(true)
		pay := New("paypal")
		_, err := per.Initiate(pay)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, SaveError) {
			t.Errorf("Expected error to be %v, got %v", SaveError, err)
		}
		spyPersistence.failWhenSavingAPayment(false)
	})
}

type spyMethod struct {
	validations    map[MethodName]int
	rejectPayments bool
	payments       map[string]Payment
	failPayments   bool
}

func (m *spyMethod) Validate(pay Payment) error {
	if m.rejectPayments {
		return errors.New("payment rejected")
	}
	m.validations[pay.Method()]++
	return nil
}

func (m *spyMethod) Create(pay Payment) (string, error) {
	if m.failPayments {
		return "", errors.New("payment failed")
	}
	id := uuid.NewString()
	m.payments[id] = pay
	return id, nil
}

func (m *spyMethod) wasValidated(pay Payment) bool {
	validations, ok := m.validations[pay.Method()]
	if !ok {
		return false
	}
	return validations > 0
}

func (m *spyMethod) rejectAllPayments(b bool) {
	m.rejectPayments = b
}

func (m *spyMethod) wasCreated(id string) bool {
	_, ok := m.payments[id]
	return ok
}

func (m *spyMethod) failAllPayments(b bool) {
	m.failPayments = b
}

func getSpyMethod() *spyMethod {
	return &spyMethod{
		validations: make(map[MethodName]int),
		payments:    make(map[string]Payment),
	}
}

func testSupportedMethod(t *testing.T, per Performer) {
	t.Helper()
	pay := New("paypal")
	_, err := per.Initiate(pay)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func testUnsupportedMethod(t *testing.T, per Performer) {
	t.Helper()
	pay := New("unsupported-method")
	_, err := per.Initiate(pay)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if !errors.Is(err, UnsupportedMethodError) {
		t.Errorf("Expected error to be %v, got %v", UnsupportedMethodError, err)
	}
}

func testPaymentWithMethod(t *testing.T, per Performer) {
	t.Helper()

	pay := New("paypal")
	_, err := per.Initiate(pay)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func testPaymentWithoutMethod(t *testing.T, per Performer) {
	t.Helper()
	var pay Payment = payment{}
	_, err := per.Initiate(pay)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if !errors.Is(err, EmptyMethodError) {
		t.Errorf("Expected error to be %v, got %v", EmptyMethodError, err)
	}
}

func testNonNilPayment(t *testing.T, per Performer) {
	t.Helper()
	pay := New("paypal")
	_, err := per.Initiate(pay)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func testNilPayment(t *testing.T, per Performer) {
	t.Helper()
	_, err := per.Initiate(nil)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if !errors.Is(err, IsNilError) {
		t.Errorf("Expected error to be %v, got %v", IsNilError, err)
	}
}

type spyPerformerPersistence struct {
	savedMethods            map[MethodName]Method
	savedPayments           map[string]Payment
	failOnPaymentSave       bool
	failOnPaymentMethodSave bool
}

func (s *spyPerformerPersistence) RetrievePaymentMethod(name MethodName) (Method, error) {
	method, ok := s.savedMethods[name]
	if !ok {
		return nil, errors.New("method not found")
	}
	return method, nil
}

func (s *spyPerformerPersistence) SavePaymentMethod(name MethodName, method Method) error {
	if s.failOnPaymentMethodSave {
		return errors.New("failed to save payment method")
	}
	s.savedMethods[name] = method
	return nil
}

func (s *spyPerformerPersistence) paymentMethodExists(name MethodName) bool {
	_, ok := s.savedMethods[name]
	return ok
}

func (s *spyPerformerPersistence) SavePayment(id string, pay Payment) error {
	if s.failOnPaymentSave {
		return errors.New("failed to save payment")
	}
	s.savedPayments[id] = pay
	return nil
}

func (s *spyPerformerPersistence) paymentExists(id string) bool {
	_, ok := s.savedPayments[id]
	return ok
}

func (s *spyPerformerPersistence) failWhenSavingAPayment(b bool) {
	s.failOnPaymentSave = b
}

func (s *spyPerformerPersistence) failWhenSavingAPaymentMethod(b bool) {
	s.failOnPaymentMethodSave = b
}

func getSpyPerformerPersistence() *spyPerformerPersistence {
	return &spyPerformerPersistence{
		savedMethods:  make(map[MethodName]Method),
		savedPayments: make(map[string]Payment),
	}
}
