package testutil

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientMock_Do_Success(t *testing.T) {
	t.Parallel()

	mockClient := MockHTTPClient(func(m *HTTPClientMock) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"success":true}`)),
		}
		m.On("Do", mock.Anything).Return(resp, nil)
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	resp, err := mockClient.Do(req)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	mockClient.AssertExpectations(t)
}

func TestHTTPClientMock_Do_Error(t *testing.T) {
	t.Parallel()

	mockClient := MockHTTPClient(func(m *HTTPClientMock) {
		m.On("Do", mock.Anything).Return(nil, errors.New("network error"))
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	resp, err := mockClient.Do(req)

	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, "network error", err.Error())

	mockClient.AssertExpectations(t)
}

func TestReadCloserMock_Read_Success(t *testing.T) {
	t.Parallel()

	mockReadCloser := MockReadCloser(func(m *ReadCloserMock) {
		m.On("Read", mock.Anything).Return(5, nil)
	})

	buf := make([]byte, 10)
	n, err := mockReadCloser.Read(buf)

	require.NoError(t, err)
	require.Equal(t, 5, n)

	mockReadCloser.AssertExpectations(t)
}

func TestReadCloserMock_Read_Error(t *testing.T) {
	t.Parallel()

	mockReadCloser := MockReadCloser(func(m *ReadCloserMock) {
		m.On("Read", mock.Anything).Return(0, io.ErrUnexpectedEOF)
	})

	buf := make([]byte, 10)
	n, err := mockReadCloser.Read(buf)

	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.Equal(t, 0, n)

	mockReadCloser.AssertExpectations(t)
}

func TestReadCloserMock_Close(t *testing.T) {
	t.Parallel()

	mockReadCloser := MockReadCloser(func(m *ReadCloserMock) {
		m.On("Close").Return(nil)
	})

	err := mockReadCloser.Close()

	require.NoError(t, err)

	mockReadCloser.AssertExpectations(t)
}

func TestReadCloserMock_Close_Error(t *testing.T) {
	t.Parallel()

	mockReadCloser := MockReadCloser(func(m *ReadCloserMock) {
		m.On("Close").Return(errors.New("close error"))
	})

	err := mockReadCloser.Close()

	require.Error(t, err)
	require.Equal(t, "close error", err.Error())

	mockReadCloser.AssertExpectations(t)
}
