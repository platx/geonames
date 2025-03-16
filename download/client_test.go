package download

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NewClient(t *testing.T) {
	t.Parallel()

	t.Run("without options", func(t *testing.T) {
		t.Parallel()

		client := NewClient()

		require.NotNil(t, client)

		assert.Equal(t, "https://download.geonames.org/export/dump", client.baseURL)

		require.IsType(t, &http.Client{}, client.httpClient)
		assert.Same(t, http.DefaultTransport, client.httpClient.(*http.Client).Transport)
		assert.Equal(t, time.Minute*10, client.httpClient.(*http.Client).Timeout)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()

		httpClient := &http.Client{}
		customURL := "http://example.com/?"

		client := NewClient(
			WithBaseURL(customURL),
			WithHTTPClient(httpClient),
			WithLogger(nopLogger{}),
		)

		require.NotNil(t, client)

		assert.Equal(t, "http://example.com", client.baseURL)
		assert.Same(t, httpClient, client.httpClient)
	})
}

type args struct {
	httpClient httpDoer
	ctx        context.Context
}

type exp[T any] struct {
	res []T
	err error
}

type testSuite[T any] struct {
	name string
	args args
	exp  exp[T]
}

func (ts testSuite[T]) run(
	t *testing.T,
	caller func(client *Client, ctx context.Context) ([]T, error),
) {
	t.Helper()

	defer mock.AssertExpectationsForObjects(t, ts.args.httpClient)

	client := NewClient(
		WithHTTPClient(ts.args.httpClient),
	)

	res, err := caller(client, ts.args.ctx)
	if ts.exp.err != nil {
		require.EqualError(t, err, ts.exp.err.Error())
		assert.Equal(t, ts.exp.res, res)

		return
	}

	require.NoError(t, err)
	assert.Equal(t, ts.exp.res, res)
}

func assertRequest(t *testing.T, req *http.Request, fileName string) bool {
	t.Helper()

	if !assert.Equal(t, http.MethodGet, req.Method) {
		return false
	}

	if !assert.Equal(t, "https", req.URL.Scheme) {
		return false
	}

	if !assert.Equal(t, "download.geonames.org", req.URL.Host) {
		return false
	}

	if !assert.Equal(t, "/export/dump/"+fileName, req.URL.Path) {
		return false
	}

	if !assert.Empty(t, req.URL.Query()) {
		return false
	}

	return true
}
