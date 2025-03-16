package webservice

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/platx/geonames/value"
)

func Test_ResponseError_Code(t *testing.T) {
	t.Parallel()

	code := value.ErrCode(42)
	err := &ResponseError{code: code}

	assert.Equal(t, code, err.Code())
}

func Test_ResponseError_Message(t *testing.T) {
	t.Parallel()

	message := "test message"
	err := &ResponseError{message: message}

	assert.Equal(t, message, err.Message())
}

func Test_ResponseError_Error(t *testing.T) {
	t.Parallel()

	err := &ResponseError{code: value.ErrCode(42), message: "test message"}

	assert.Equal(t, "got error response => code: 42, message: \"test message\"", err.Error())
}

func Test_ResponseError_MatchCode(t *testing.T) {
	t.Parallel()

	code := value.ErrCode(42)
	err := &ResponseError{code: code}

	assert.True(t, err.MatchCode(code))
	assert.False(t, err.MatchCode(value.ErrCode(43)))
}

func Test_error_As_ResponseError(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		source := &ResponseError{}

		var target *ResponseError

		assert.ErrorAs(t, source, &target)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		source := errors.New("test error")

		var target *ResponseError

		assert.False(t, errors.As(source, &target))
	})
}
