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

func Test_Client_CountrySubdivision(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req CountrySubdivisionRequest) (CountrySubdivision, error) {
		return client.CountrySubdivision
	}

	testCases := []testSuite[CountrySubdivisionRequest, CountrySubdivision]{
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
								"/countrySubdivisionJSON",
								url.Values{
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"radius":   []string{"11"},
									"level":    []string{"1"},
									"lang":     []string{"en"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "country_subdivision.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountrySubdivisionRequest]{
				ctx: context.Background(),
				req: CountrySubdivisionRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Radius:   11,
					Level:    1,
					Language: "en",
				},
			},
			exp: exp[CountrySubdivision]{
				res: CountrySubdivision{
					GeoNameID: 1,
					Country: value.Country{
						Code: value.CountryCodeUnitedStates,
						Name: "United States",
					},
					Codes: []value.AdminLevelCode{
						{
							Level: 1,
							Type:  "XXXX",
							Code:  "XX",
						},
						{
							Level: 2,
							Type:  "XXXY",
							Code:  "XY",
						},
					},
					AdminDivision: value.AdminDivisions{
						First: value.AdminDivision{
							ID:   11,
							Code: "D1",
							Name: "Test division 11",
						},
						Second: value.AdminDivision{
							ID:   12,
							Code: "D2",
							Name: "Test division 12",
						},
						Third: value.AdminDivision{
							ID:   13,
							Code: "D3",
							Name: "Test division 13",
						},
						Fourth: value.AdminDivision{
							ID:   14,
							Code: "D4",
							Name: "Test division 14",
						},
						Fifth: value.AdminDivision{
							ID:   15,
							Code: "D5",
							Name: "Test division 15",
						},
					},
					Distance: 1.11,
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
