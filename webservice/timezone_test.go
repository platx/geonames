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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_Timezone(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req TimezoneRequest) (Timezone, error) {
		return client.Timezone
	}

	testCases := []testSuite[TimezoneRequest, Timezone]{
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
								"/timezoneJSON",
								url.Values{
									"lat":      []string{"1.11"},
									"lng":      []string{"-1.11"},
									"radius":   []string{"11"},
									"lang":     []string{"en"},
									"date":     []string{"2021-01-01"},
									"type":     []string{"json"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "timezone.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[TimezoneRequest]{
				ctx: context.Background(),
				req: TimezoneRequest{
					Position: value.Position{
						Latitude:  1.11,
						Longitude: -1.11,
					},
					Radius:   11,
					Language: "en",
					Date:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			exp: exp[Timezone]{
				res: Timezone{
					Name: "UTC",
					Country: value.Country{
						Code: value.CountryCodeUnitedKingdom,
						Name: "United Kingdom",
					},
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Time:      time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
					Sunset:    time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
					Sunrise:   time.Date(2021, 1, 1, 22, 0, 0, 0, time.UTC),
					GMTOffset: 1,
					DSTOffset: 2,
					RawOffset: 3,
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
						Body:       io.NopCloser(strings.NewReader(`{"`)),
					})
				}),
				userName: "test-user",
			},
			args: args[TimezoneRequest]{
				ctx: context.Background(),
				req: TimezoneRequest{},
			},
			exp: exp[Timezone]{
				res: Timezone{},
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
			args: args[TimezoneRequest]{
				ctx: context.Background(),
				req: TimezoneRequest{},
			},
			exp: exp[Timezone]{
				res: Timezone{},
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
			args: args[TimezoneRequest]{
				ctx: context.Background(),
				req: TimezoneRequest{},
			},
			exp: exp[Timezone]{
				res: Timezone{},
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
			args: args[TimezoneRequest]{
				ctx: context.Background(),
				req: TimezoneRequest{},
			},
			exp: exp[Timezone]{
				res: Timezone{},
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[TimezoneRequest]{
				ctx: nil,
				req: TimezoneRequest{},
			},
			exp: exp[Timezone]{
				res: Timezone{},
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
