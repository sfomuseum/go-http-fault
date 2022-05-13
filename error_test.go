package fault

import (
	"fmt"
	"testing"
)

type ErrorClass string

const ExpiredCode ErrorClass = "Code has expired."

const UnknownCode ErrorClass = "Unknown code."

type TestError struct {
	FaultError
	class ErrorClass
	error error
}

func (e *TestError) Public() error {
	return fmt.Errorf("%s", e.class)
}

func (e *TestError) Private() error {
	return e.error
}

func NewTestError(cl ErrorClass, err error) FaultError {

	e := &TestError{
		class: cl,
		error: err,
	}

	return e
}

func TestFaultError(t *testing.T) {

	e := NewTestError(UnknownCode, fmt.Errorf("SNFU"))

	pub := e.Public()
	pri := e.Private()

	if pub.Error() != string(UnknownCode) {
		t.Fatalf("Invalid public error: %s", pub)
	}

	if pri.Error() != "SNFU" {
		t.Fatalf("Invalid private error: %s", pub)
	}

}
