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

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_Astergdem(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, position value.Position) (int32, error) {
		return client.Astergdem
	}

	testCases := []testSuite[value.Position, int32]{
		{
			name: "success with request values",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/astergdemJSON",
								url.Values{
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"type":     []string{"json"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "astergdem.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[value.Position]{
				ctx: context.Background(),
				req: value.Position{
					Latitude:  1.111,
					Longitude: -1.111,
				},
			},
			exp: exp[int32]{
				res: 111,
				err: nil,
			},
		},
		{
			name: "invalid success response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"`)),
					})
				}),
				userName: "test-user",
			},
			args: args[value.Position]{
				ctx: context.Background(),
				req: value.Position{},
			},
			exp: exp[int32]{
				res: 0,
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
			args: args[value.Position]{
				ctx: context.Background(),
				req: value.Position{},
			},
			exp: exp[int32]{
				res: 0,
				err: errors.New("decode response => got error response => code: 10, message: \"user does not exist.\""),
			},
		},
		{
			name: "invalid error response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(strings.NewReader(`{"`)),
					})
				}),
				userName: "test-user",
			},
			args: args[value.Position]{
				ctx: context.Background(),
				req: value.Position{},
			},
			exp: exp[int32]{
				res: 0,
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
			args: args[value.Position]{
				ctx: context.Background(),
				req: value.Position{},
			},
			exp: exp[int32]{
				res: 0,
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[value.Position]{
				ctx: nil,
				req: value.Position{},
			},
			exp: exp[int32]{
				res: 0,
				err: errors.New("create http request => net/http: nil Context"),
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
