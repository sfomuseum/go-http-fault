package fault

import (
	"fmt"
	"html/template"
	"io"
	"log"
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

func TestFaultHandlerWrapper(t *testing.T) {

	logger := log.Default()

	fault_t, err := template.New("fault").Parse(`{{ .Status }} {{ .Error }}`)

	if err != nil {
		t.Fatalf("Failed to parse template, %v", err)
	}

	fw := NewFaultWrapper(logger, fault_t)

	mux := http.NewServeMux()

	ok_h := okHandler()
	err_h := errorHandler()

	fw.HandleWithMux(mux, "/ok", ok_h)
	fw.HandleWithMux(mux, "/err", err_h)

	go func() {

		err := http.ListenAndServe("localhost:8080", mux)

		if err != nil {
			t.Fatalf("Failed to serve requests, %v", err)
		}
	}()

	req, err := http.Get("http://localhost:8080/ok")

	if err != nil {
		t.Fatalf("Failed to GET /ok, %v", err)
	}

	if req.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected response for ok handler, %v", req.StatusCode)
	}

	req2, err := http.Get("http://localhost:8080/err")

	if err != nil {
		t.Fatalf("Failed to GET /err, %v", err)
	}

	body, err := io.ReadAll(req2.Body)

	if err != nil {
		t.Fatalf("Failed to read response for /err, %v", err)
	}

	if string(body) != "999 SAD FACE" {
		t.Fatalf("Unexpected response, '%s'", string(body))
	}

}
