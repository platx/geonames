package webservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/webservice/testdata"
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

func Test_Client_apiRequest(t *testing.T) {
	t.Parallel()

	type testRequest struct {
		Key1 string `url:"key1"`
		Key2 int    `url:"key2"`
	}

	type testResult struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	caller := func(client *Client) func(ctx context.Context, req testRequest) (testResult, error) {
		return func(ctx context.Context, req testRequest) (testResult, error) {
			var res testResult

			err := client.apiRequest(
				ctx,
				"/test",
				req,
				&res,
			)

			return res, err
		}
	}

	testCases := []testSuite[testRequest, testResult]{
		{
			name: "success",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/test",
								url.Values{
									"key1":     []string{"val1"},
									"key2":     []string{"2"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"key1":"val1","key2":2}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[testRequest]{
				ctx: context.Background(),
				req: testRequest{
					Key1: "val1",
					Key2: 2,
				},
			},
			exp: exp[testResult]{
				res: testResult{
					Key1: "val1",
					Key2: 2,
				},
				err: nil,
			},
		},
		{
			name: "invalid success response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geo`)),
					})
				}),
				userName: "test-user",
			},
			args: args[testRequest]{
				ctx: context.Background(),
				req: testRequest{},
			},
			exp: exp[testResult]{
				res: testResult{},
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "error response",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       testutil.MustOpen(testdata.FS, "authorization_error.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[testRequest]{
				ctx: context.Background(),
				req: testRequest{},
			},
			exp: exp[testResult]{
				res: testResult{},
				err: errors.New("decode response => got error response => code: 10, message: \"user does not exist.\""),
			},
		},
		{
			name: "invalid error response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(strings.NewReader(`{"stat`)),
					})
				}),
				userName: "test-user",
			},
			args: args[testRequest]{
				ctx: context.Background(),
				req: testRequest{},
			},
			exp: exp[testResult]{
				res: testResult{},
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "send request failed",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				userName: "test-user",
			},
			args: args[testRequest]{
				ctx: context.Background(),
				req: testRequest{},
			},
			exp: exp[testResult]{
				res: testResult{},
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[testRequest]{
				ctx: nil,
				req: testRequest{},
			},
			exp: exp[testResult]{
				res: testResult{},
				err: errors.New("create http request => net/http: nil Context"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
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
