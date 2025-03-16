package webservice

import (
	"fmt"

	"github.com/platx/geonames/value"
)

type ResponseError struct {
	code    value.ErrCode
	message string
}

func (e *ResponseError) Code() value.ErrCode {
	return e.code
}

func (e *ResponseError) Message() string {
	return e.message
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("got error response => code: %d, message: %q", e.Code(), e.Message())
}

func (e *ResponseError) MatchCode(code value.ErrCode) bool {
	return e.code == code
}
