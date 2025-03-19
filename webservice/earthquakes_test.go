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

func Test_Client_Earthquakes(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req EarthquakesRequest) ([]Earthquake, error) {
		return client.Earthquakes
	}

	testCases := []testSuite[EarthquakesRequest, []Earthquake]{
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
								"/earthquakesJSON",
								url.Values{
									"east":         []string{"1.1"},
									"west":         []string{"1.2"},
									"north":        []string{"-1.1"},
									"south":        []string{"-1.2"},
									"date":         []string{"2021-01-02"},
									"minMagnitude": []string{"1.1"},
									"maxRows":      []string{"2"},
									"type":         []string{"json"},
									"username":     []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "earthquakes.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{
					BoundingBox: value.BoundingBox{
						East:  1.1,
						West:  1.2,
						North: -1.1,
						South: -1.2,
					},
					Date:         time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					MinMagnitude: 1.1,
					MaxRows:      2,
				},
			},
			exp: exp[[]Earthquake]{
				res: []Earthquake{
					{
						ID: "x0001yyy",
						Position: value.Position{
							Latitude:  1.111,
							Longitude: -1.111,
						},
						Depth:     11.1,
						Source:    "xx",
						Magnitude: 1.1,
						Time:      time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
					},
					{
						ID: "x0002yyy",
						Position: value.Position{
							Latitude:  2.222,
							Longitude: -2.222,
						},
						Depth:     22.2,
						Source:    "yy",
						Magnitude: 2.2,
						Time:      time.Date(2021, 1, 2, 12, 0, 0, 0, time.UTC),
					},
				},
				err: nil,
			},
		},
		{
			name: "empty without request values",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/earthquakesJSON",
								url.Values{
									"type":     []string{"json"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "earthquakes_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: []Earthquake{},
				err: nil,
			},
		},
		{
			name: "invalid time",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"earthquakes": [{"datetime": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: nil,
				err: errors.New("decode response => parse Time => parsing time \"invalid\" as \"2006-01-02 15:04:05\": cannot parse \"invalid\" as \"2006\""),
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
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: nil,
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
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: nil,
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
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: nil,
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
			args: args[EarthquakesRequest]{
				ctx: context.Background(),
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: nil,
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[EarthquakesRequest]{
				ctx: nil,
				req: EarthquakesRequest{},
			},
			exp: exp[[]Earthquake]{
				res: nil,
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
