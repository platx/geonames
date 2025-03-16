package webservice

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NewClient(t *testing.T) {
	t.Parallel()

	t.Run("without options", func(t *testing.T) {
		t.Parallel()

		userName := "test-user"

		client := NewClient(userName)

		require.NotNil(t, client)

		assert.Equal(t, "https://secure.geonames.org", client.baseURL)
		assert.Equal(t, userName, client.userName)

		require.IsType(t, &http.Client{}, client.httpClient)
		assert.Same(t, http.DefaultTransport, client.httpClient.(*http.Client).Transport)
		assert.Equal(t, defaultRequestTimeout, client.httpClient.(*http.Client).Timeout)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()

		httpClient := &http.Client{}
		userName := "test-user"
		customURL := "http://api.geonames.org/?"

		client := NewClient(
			userName,
			WithBaseURL(customURL),
			WithHTTPClient(httpClient),
		)

		require.NotNil(t, client)

		assert.Equal(t, "http://api.geonames.org", client.baseURL)
		assert.Equal(t, userName, client.userName)
		assert.Same(t, httpClient, client.httpClient)
	})
}

type deps struct {
	httpClient httpDoer
	userName   string
}

type args[T any] struct {
	ctx context.Context
	req T
}

type exp[T any] struct {
	res T
	err error
}

type testSuite[REQ, RES any] struct {
	name string
	deps deps
	args args[REQ]
	exp  exp[RES]
}

func (ts testSuite[REQ, RES]) run(
	t *testing.T,
	caller func(client *Client) func(ctx context.Context, req REQ) (RES, error),
) {
	t.Helper()

	defer mock.AssertExpectationsForObjects(t, ts.deps.httpClient)

	client := NewClient(
		ts.deps.userName,
		WithHTTPClient(ts.deps.httpClient),
	)

	res, err := caller(client)(ts.args.ctx, ts.args.req)
	if ts.exp.err != nil {
		require.EqualError(t, err, ts.exp.err.Error())
		assert.Equal(t, ts.exp.res, res)

		return
	}

	require.NoError(t, err)
	assert.Equal(t, ts.exp.res, res)
}

func assertRequest(t *testing.T, req *http.Request, path string, urlValues url.Values) bool {
	t.Helper()

	if !assert.Equal(t, http.MethodGet, req.Method) {
		return false
	}

	if !assert.Equal(t, "https", req.URL.Scheme) {
		return false
	}

	if !assert.Equal(t, "secure.geonames.org", req.URL.Host) {
		return false
	}

	if !assert.Equal(t, path, req.URL.Path) {
		return false
	}

	if !assert.Equal(t, urlValues, req.URL.Query()) {
		return false
	}

	return true
}
