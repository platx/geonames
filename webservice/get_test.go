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

func Test_Client_Get(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req GetRequest) (GeoNameDetailed, error) {
		return client.Get
	}

	testCases := []testSuite[GetRequest, GeoNameDetailed]{
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
								"/getJSON",
								url.Values{
									"geonameId": []string{"1"},
									"lang":      []string{"en"},
									"username":  []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geoname_detailed.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[GetRequest]{
				ctx: context.Background(),
				req: GetRequest{
					ID:       1,
					Language: "en",
				},
			},
			exp: exp[GeoNameDetailed]{
				res: GeoNameDetailed{
					GeoName: GeoName{
						ID: 1,
						Country: value.Country{
							Code: value.CountryCodeUnitedKingdom,
							Name: "United Kingdom",
						},
						AdminSubdivision: value.AdminDivisions{
							First: value.AdminDivision{
								Code: "FOO",
								Name: "Foo",
							},
							Second: value.AdminDivision{
								Code: "BAR",
								Name: "Bar",
							},
							Third: value.AdminDivision{
								Code: "BAZ",
								Name: "Baz",
							},
							Fourth: value.AdminDivision{
								Code: "FOOBAR",
								Name: "FooBar",
							},
							Fifth: value.AdminDivision{},
						},
						Feature: value.Feature{
							Class:     "A",
							ClassName: "Test class",
							Code:      "AAAA",
							CodeName:  "Test code",
						},
						Name:        "London",
						ToponymName: "London",
						Position: value.Position{
							Latitude:  1.111,
							Longitude: -1.111,
						},
						Population: 111111,
					},
					ContinentCode: value.ContinentCodeEurope,
					ASCIIName:     "London",
					AlternateNames: []value.AlternateName{
						{
							Language: "link",
							Value:    "https://example.com/london",
						},
						{
							Language: "om",
							Value:    "Landan",
						},
						{
							Language: "en",
							Value:    "London",
						},
						{
							Language: "ru",
							Value:    "Лондон",
						},
					},
					Timezone: value.Timezone{
						Name:      "Europe/London",
						GMTOffset: 1,
						DSTOffset: 2,
					},
					Elevation: 111,
					SRTM3:     11,
					Astergdem: 12,
					BoundingBox: value.BoundingBox{
						East:  1.1,
						West:  1.2,
						North: -1.1,
						South: -1.2,
					},
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
