package fault

import (
	"fmt"
	"net/http"
	"testing"
)

func okHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		return
	}

	return http.HandlerFunc(fn)
}

func errorHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		err := fmt.Errorf("SAD FACE")
		code := 999

		AssignError(req, err, code)
		rsp.WriteHeader(http.StatusBadRequest)

		return
	}

	return http.HandlerFunc(fn)
}

func TestAssignError(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("Failed to create new request, %v", err)
	}

	err = fmt.Errorf("Testing")
	code := 999

	AssignError(req, err, code)

	code2, err2 := RetrieveError(req)

	if code2 != code {
		t.Fatalf("Invalid status code returned. Expected %d but got %d", code, code2)
	}

	if err2.Error() != err.Error() {
		t.Fatalf("Invalid error returned. Expected '%s' but got '%s'", err.Error(), err2.Error())
	}
}

func TestFaultHandler(t *testing.T) {
	t.Skip()
}
