package testutil

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HTTPClientMock struct {
	mock.Mock
}

func (m *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	if res := args.Get(0); res != nil {
		return res.(*http.Response), nil
	}

	return nil, args.Error(1)
}

func MockHTTPClient(fn func(m *HTTPClientMock)) *HTTPClientMock {
	m := &HTTPClientMock{}

	fn(m)

	return m
}

type ReadCloserMock struct {
	mock.Mock
}

func (m *ReadCloserMock) Read(p []byte) (int, error) {
	args := m.Called(p)

	return args.Int(0), args.Error(1)
}

func (m *ReadCloserMock) Close() error {
	return m.Called().Error(0)
}

func MockReadCloser(fn func(m *ReadCloserMock)) *ReadCloserMock {
	m := &ReadCloserMock{}

	fn(m)

	return m
}
