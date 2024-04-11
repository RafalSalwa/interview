// Package path implements HTTP responses struct features and functions
package responses

import (
	"log"
	"net/http"
)

type Builder struct {
	Response       ErrResponse
	ResponseWriter http.ResponseWriter
}

type ErrResponse struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Headers map[string]interface{}
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    data   `json:"data"`
	Err     error
}

func NewErrorBuilder() *Builder {
	return &Builder{
		Response: ErrResponse{},
	}
}

func (rb *Builder) SetResponseCode(statusCode int) *Builder {
	rb.Response.Code = statusCode
	return rb
}

func (rb *Builder) SetReason(reason string) *Builder {
	rb.Response.Reason = reason
	return rb
}

func (rb *Builder) AddHeader(key, val string) *Builder {
	rb.Response.Headers[key] = val
	return rb
}

func (rb *Builder) SetMessage(msg string) *Builder {
	rb.Response.Message = msg
	return rb
}

func (rb *Builder) SetWriter(w http.ResponseWriter) *Builder {
	rb.ResponseWriter = w
	return rb
}

func (rb *Builder) SetError(err error) *Builder {
	rb.Response.Err = err
	return rb
}

func (rb *Builder) Respond() {
	body := marshalErrorResponse(rb.Response.Reason)

	rb.ResponseWriter.WriteHeader(rb.Response.Code)
	_, err := rb.ResponseWriter.Write(body)
	if err != nil {
		rb.Response.Err = err
		log.Fatal(err)
	}
}
