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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}
