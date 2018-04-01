package result

import (
	"fmt"
	"io"
)

type (
	Result struct {
		Body  ResultBody
		Error error
	}
	ResultBody interface {
		JSON() interface{}
		Pretty(writer io.Writer) error
		Plain(writer io.Writer) error
	}
)

func New(body ResultBody) *Result {
	return &Result{body, nil}
}

func NewError(err error) *Result {
	return &Result{nil, err}
}

func (r *Result) JSON() map[string]interface{} {
	if r == nil {
		return map[string]interface{}{
			"err": "result not set.  this is an implementation problem.",
		}
	}
	v := map[string]interface{}{}
	if r.Error != nil {
		v["err"] = r.Error.Error()
	}
	if r.Body != nil {
		body := r.Body.JSON()
		if body != nil {
			v["body"] = r.Body.JSON()
		}
	}
	return v
}

func (r *Result) Pretty(writer io.Writer) {
	if r.Error == nil {
		r.Body.Pretty(writer)
	} else {
		fmt.Fprintln(writer, r.Error)
	}
}

func (r *Result) Plain(writer io.Writer) {
	if r.Error == nil {
		r.Body.Plain(writer)
	} else {
		fmt.Fprintln(writer, r.Error)
	}
}
