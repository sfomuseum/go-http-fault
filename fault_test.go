package fault

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"testing"
)

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

func TestFaultHandlerVarsFunc(t *testing.T) {

	type ValidCustomVars struct {
		Custom string
		FaultHandlerVars
	}

	type InvalidCustomVars struct {
		Custom string
	}

	valid_func := func() interface{} {

		vars := &ValidCustomVars{
			Custom: "custom",
		}

		return vars
	}

	invalid_func := func() interface{} {

		vars := &InvalidCustomVars{
			Custom: "custom",
		}

		return vars
	}

	valid_vars := valid_func()

	if !ImplementsFaultHandlerVars(valid_vars) {
		t.Fatalf("%T does not implement fault handler vars", valid_vars)
	}

	invalid_vars := invalid_func()

	if ImplementsFaultHandlerVars(invalid_vars) {
		t.Fatalf("%T implements fault handler vars, which is not expected", invalid_vars)
	}
}

func TestFaultHandlerWithCustomVars(t *testing.T) {

	tpl := template.New("test")
	tpl, err := tpl.Parse(`{{ .Custom }} {{ .Status }}`)

	if err != nil {
		t.Fatalf("Failed to parse template, %v", err)
	}

	type CustomVars struct {
		Custom string
		FaultHandlerVars
	}

	custom_func := func() interface{} {

		vars := &CustomVars{
			Custom: "This is custom text",
		}

		return vars
	}

	opts := &FaultHandlerOptions{
		Logger:   log.Default(),
		Template: tpl,
		VarsFunc: custom_func,
	}

	fh := FaultHandlerWithOptions(opts)

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		err := fmt.Errorf("This is a test")
		AssignError(req, err, 999)
		fh.ServeHTTP(rsp, req)
		return
	}

	h := http.HandlerFunc(fn)

	mux := http.NewServeMux()
	mux.Handle("/", h)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	defer s.Close()

	go func() {
		s.ListenAndServe()
	}()

	rsp, err := http.Get("http://localhost:8080")

	if err != nil {
		t.Fatalf("Failed to get from localhost, %v", err)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Fatalf("Failed to read response, %v", err)
	}

	str_body := string(body)
	expected_body := "This is custom text 999"

	if str_body != expected_body {
		t.Fatalf("Unexpected output '%s' (got '%s')", str_body, expected_body)
	}

}
