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

func Test_Client_WeatherICAO(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req WeatherICAORequest) (WeatherObservation, error) {
		return client.WeatherICAO
	}

	testCases := []testSuite[WeatherICAORequest, WeatherObservation]{
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
								"/weatherIcaoJSON",
								url.Values{
									"ICAO":     []string{"XXXX"},
									"lang":     []string{"en"},
									"type":     []string{"json"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "weather_icao.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[WeatherICAORequest]{
				ctx: context.Background(),
				req: WeatherICAORequest{
					ICAO:     "XXXX",
					Language: "en",
				},
			},
			exp: exp[WeatherObservation]{
				res: WeatherObservation{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Observation:      "Test observation 1",
					ICAO:             "XXXX",
					StationName:      "Test station 1",
					CloudsCode:       "XX",
					CloudsName:       "test clouds 1",
					WeatherCondition: "test condition 1",
					Temperature:      1,
					DewPoint:         -1,
					Humidity:         11,
					WindDirection:    111,
					WindSpeed:        1,
					UpdatedAt:        time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
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
			args: args[WeatherICAORequest]{
				ctx: context.Background(),
				req: WeatherICAORequest{},
			},
			exp: exp[WeatherObservation]{
				res: WeatherObservation{},
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
			args: args[WeatherICAORequest]{
				ctx: context.Background(),
				req: WeatherICAORequest{},
			},
			exp: exp[WeatherObservation]{
				res: WeatherObservation{},
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
			args: args[WeatherICAORequest]{
				ctx: context.Background(),
				req: WeatherICAORequest{},
			},
			exp: exp[WeatherObservation]{
				res: WeatherObservation{},
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
			args: args[WeatherICAORequest]{
				ctx: context.Background(),
				req: WeatherICAORequest{},
			},
			exp: exp[WeatherObservation]{
				res: WeatherObservation{},
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[WeatherICAORequest]{
				ctx: nil,
				req: WeatherICAORequest{},
			},
			exp: exp[WeatherObservation]{
				res: WeatherObservation{},
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
