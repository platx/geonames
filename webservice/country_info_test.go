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

func Test_Client_CountryInfo(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req CountryInfoRequest) ([]CountryDetailed, error) {
		return client.CountryInfo
	}

	testCases := []testSuite[CountryInfoRequest, []CountryDetailed]{
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
								"/countryInfoJSON",
								url.Values{
									"country":  []string{"DE", "UA"},
									"lang":     []string{"en"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "country_detailed.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{
					Country:  []value.CountryCode{value.CountryCodeGermany, value.CountryCodeUkraine},
					Language: "en",
				},
			},
			exp: exp[[]CountryDetailed]{
				res: []CountryDetailed{
					{
						Country: value.Country{
							ID:   1,
							Code: value.CountryCodeGermany,
							Name: "Germany",
						},
						Continent: value.Continent{
							Code: value.ContinentCodeEurope,
							Name: "Europe",
						},
						Capital:   "Berlin",
						Languages: []string{"de"},
						BoundingBox: value.BoundingBox{
							East:  1.1,
							West:  1.2,
							North: -1.1,
							South: -1.2,
						},
						IsoAlpha3:        "DEU",
						IsoNumeric:       276,
						FipsCode:         "GM",
						Population:       111111,
						AreaInSqKm:       11111.1,
						PostalCodeFormat: "XXXXX",
						CurrencyCode:     "EUR",
					},
					{
						Country: value.Country{
							ID:   2,
							Code: value.CountryCodeUkraine,
							Name: "Ukraine",
						},
						Continent: value.Continent{
							Code: value.ContinentCodeEurope,
							Name: "Europe",
						},
						Capital:   "Kyiv",
						Languages: []string{"uk", "ru-UA"},
						BoundingBox: value.BoundingBox{
							East:  2.1,
							West:  2.2,
							North: -2.1,
							South: -2.2,
						},
						IsoAlpha3:        "UKR",
						IsoNumeric:       804,
						FipsCode:         "UP",
						Population:       222222,
						AreaInSqKm:       22222.2,
						PostalCodeFormat: "XXXXX",
						CurrencyCode:     "UAH",
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
								"/countryInfoJSON",
								url.Values{
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geonames_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: []CountryDetailed{},
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
