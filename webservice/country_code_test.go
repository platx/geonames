package webservice

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_CountryCode(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req CountryCodeRequest) (CountryNearby, error) {
		return client.CountryCode
	}

	testCases := []testSuite[CountryCodeRequest, CountryNearby]{
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
								"/countryCodeJSON",
								url.Values{
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"radius":   []string{"10"},
									"lang":     []string{"en"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "country_nearby.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryCodeRequest]{
				ctx: context.Background(),
				req: CountryCodeRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Radius:   10,
					Language: "en",
				},
			},
			exp: exp[CountryNearby]{
				res: CountryNearby{
					Country: value.Country{
						Code: value.CountryCodeUnitedStates,
						Name: "United States",
					},
					Distance:  0.111,
					Languages: []string{"en", "es"},
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
								"/countryCodeJSON",
								url.Values{
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "country_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryCodeRequest]{
				ctx: context.Background(),
				req: CountryCodeRequest{},
			},
			exp: exp[CountryNearby]{
				res: CountryNearby{Languages: []string{}},
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
