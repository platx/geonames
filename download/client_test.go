package download

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
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

func Test_Client_downloadAndParseFile(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, error) {
		res := make([]GeoName, 0)

		err := client.downloadAndParseFile(ctx, "countryInfo.txt", func(_ []string) error {
			assert.Fail(t, "should not be called")

			return nil
		})

		return res, err
	}

	testCases := []testSuite[GeoName]{
		{
			name: "context canceled",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "countryInfo.txt"),
						},
						nil,
					)
				}),
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				}(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file => context canceled"),
			},
		},
		{
			name: "content copy failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body: testutil.MockReadCloser(func(m *testutil.ReadCloserMock) {
								m.On("Read", mock.Anything).Return(0, assert.AnError)
								m.On("Close").Return(nil)
							}),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => copy file content => assert.AnError general error for testing"),
			},
		},
		{
			name: "invalid status code",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => unexpected status code: 500"),
			},
		},
		{
			name: "http send request failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => http client do => assert.AnError general error for testing"),
			},
		},
		{
			name: "missing context",
			args: args{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				ctx:        nil,
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => create http request => net/http: nil Context"),
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}

func Test_Client_downloadAndParseZIPFile(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, error) {
		res := make([]GeoName, 0)

		err := client.AllCountries(ctx, func(parsed GeoName) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCases := []testSuite[GeoName]{
		{
			name: "context canceled",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "allCountries.zip"),
						},
						nil,
					)
				}),
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				}(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => parse file => context canceled"),
			},
		},
		{
			name: "missing target file in archive",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "missing.zip"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => file not found in archive"),
			},
		},
		{
			name: "invalid archive",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "invalid.zip"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => open zip archive => zip: not a valid zip file"),
			},
		},
		{
			name: "no content",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => open zip archive => zip: not a valid zip file"),
			},
		},
		{
			name: "content copy failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body: testutil.MockReadCloser(func(m *testutil.ReadCloserMock) {
								m.On("Read", mock.Anything).Return(0, assert.AnError)
								m.On("Close").Return(nil)
							}),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => copy file content => assert.AnError general error for testing"),
			},
		},
		{
			name: "invalid status code",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => unexpected status code: 500"),
			},
		},
		{
			name: "http send request failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => http client do => assert.AnError general error for testing"),
			},
		},
		{
			name: "missing context",
			args: args{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				ctx:        nil,
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => create http request => net/http: nil Context"),
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
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
