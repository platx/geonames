package webservice

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_WeatherICAO(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req WeatherICAORequest) (WeatherObservationNearby, error) {
		return client.WeatherICAO
	}

	testCases := []testSuite[WeatherICAORequest, WeatherObservationNearby]{
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
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "weather_nearby.json"),
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
			exp: exp[WeatherObservationNearby]{
				res: WeatherObservationNearby{
					WeatherObservation: WeatherObservation{
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
					CountryCode:        value.CountryCodeUnitedStates,
					Elevation:          1111,
					HectoPascAltimeter: 11,
				},
				err: nil,
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
