package fault

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

const error_key string = "github.com/sfomuseum/go-http-fault#error"
const status_key string = "github.com/sfomuseum/go-http-fault#status"

func AssignError(req *http.Request, err error, status int) *http.Request {

	ctx := req.Context()
	ctx = context.WithValue(ctx, error_key, err)
	ctx = context.WithValue(ctx, status_key, status)
	return req.WithContext(ctx)
}

func FaultHandler(wr io.Writer) (http.Handler, error) {

	prefix := ""
	flags := log.LstdFlags

	l := log.New(wr, prefix, flags)

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		err_v := ctx.Value(error_key)
		status_v := ctx.Value(status_key)

		var msg string
		var status int

		if err_v == nil {
			msg = "FaultHandler triggered without an error context."
		} else {
			err := err_v.(error)
			msg = err.Error()
		}

		if status_v == nil {
			status = http.StatusInternalServerError
		} else {
			status = status_v.(int)
		}

		addr := req.RemoteAddr

		l.Printf("[FAULT] %s \"%s %s %s\" %s\n", addr, req.Method, req.RequestURI, req.Proto, msg)

		err_msg := fmt.Sprintf("There was a problem completing your request (%d)", status)

		http.Error(rsp, err_msg, status)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
