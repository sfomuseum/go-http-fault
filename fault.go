package fault

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

const ErrorKey string = "github.com/sfomuseum/go-http-fault#error"
const StatusKey string = "github.com/sfomuseum/go-http-fault#status"

type FaultHandlerVars struct {
	Status int
	Error  error
}

func AssignError(req *http.Request, err error, status int) *http.Request {

	ctx := req.Context()
	ctx = context.WithValue(ctx, ErrorKey, err)
	ctx = context.WithValue(ctx, StatusKey, status)
	return req.WithContext(ctx)
}

func RetrieveError(req *http.Request) (int, error) {

	ctx := req.Context()
	err_v := ctx.Value(ErrorKey)
	status_v := ctx.Value(StatusKey)

	var status int
	var err error

	if err_v == nil {
		msg := "FaultHandler triggered without an error context."
		err = errors.New(msg)
	} else {
		err = err_v.(error)
	}

	if status_v == nil {
		status = http.StatusInternalServerError
	} else {
		status = status_v.(int)
	}

	return status, err
}

func FaultHandler(wr io.Writer) (http.Handler, error) {
	return faultHandler(wr, nil)
}

func TemplatedFaultHandler(wr io.Writer, t *template.Template) (http.Handler, error) {
	return faultHandler(wr, t)
}

func faultHandler(wr io.Writer, t *template.Template) (http.Handler, error) {

	prefix := ""
	flags := log.LstdFlags

	l := log.New(wr, prefix, flags)

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		status, err := RetrieveError(req)
		msg := err.Error()

		addr := req.RemoteAddr

		l.Printf("[FAULT] %s \"%s %s %s\" %s\n", addr, req.Method, req.RequestURI, req.Proto, msg)

		if t != nil {

			rsp.Header().Set("Content-Type", "text/html")
			
			vars := FaultHandlerVars{
				Status: status,
				Error:  err,
			}

			err = t.Execute(rsp, vars)

			if err == nil {
				return
			}

			msg := fmt.Sprintf("Failed to render template for fault handler, %v", err)
			l.Printf("[FAULT] %s \"%s %s %s\" %s\n", addr, req.Method, req.RequestURI, req.Proto, msg)
		}

		err_msg := fmt.Sprintf("There was a problem completing your request (%d)", status)

		http.Error(rsp, err_msg, status)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
