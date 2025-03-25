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

func Test_Client_Neighbours(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req NeighboursRequest) ([]GeoName, error) {
		return client.Neighbours
	}

	testCases := []testSuite[NeighboursRequest, []GeoName]{
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
								"/neighboursJSON",
								url.Values{
									"geonameId": []string{"1"},
									"country":   []string{"US"},
									"username":  []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geonames.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[NeighboursRequest]{
				ctx: context.Background(),
				req: NeighboursRequest{
					ID:      1,
					Country: value.CountryCodeUnitedStates,
				},
			},
			exp: exp[[]GeoName]{
				res: []GeoName{
					{
						ID: 1,
						Country: value.Country{
							ID:   11,
							Code: value.CountryCodeUnitedStates,
							Name: "United States",
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
						Name:        "New York City",
						ToponymName: "New York City",
						Position: value.Position{
							Latitude:  1.111,
							Longitude: -1.111,
						},
						Population: 111111,
					},
					{
						ID: 2,
						Country: value.Country{
							ID:   22,
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
							Latitude:  2.222,
							Longitude: -2.222,
						},
						Population: 222222,
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
