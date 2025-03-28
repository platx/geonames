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

func Test_Client_FindNearbyPostalCodes(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req FindNearbyPostalCodesRequest) ([]PostalCodeNearby, error) {
		return client.FindNearbyPostalCodes
	}

	testCases := []testSuite[FindNearbyPostalCodesRequest, []PostalCodeNearby]{
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
								"/findNearbyPostalCodesJSON",
								url.Values{
									"lat":          []string{"1.111"},
									"lng":          []string{"-1.111"},
									"radius":       []string{"10"},
									"maxRows":      []string{"2"},
									"country":      []string{"US"},
									"localCountry": []string{"true"},
									"username":     []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "postalcodes_nearby.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyPostalCodesRequest]{
				ctx: context.Background(),
				req: FindNearbyPostalCodesRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Radius:       10,
					MaxRows:      2,
					Country:      value.CountryCodeUnitedStates,
					LocalCountry: true,
				},
			},
			exp: exp[[]PostalCodeNearby]{
				res: []PostalCodeNearby{
					{
						PostalCode: PostalCode{
							Code:        "1111",
							CountryCode: "US",
							AdminDivisions: value.AdminDivisions{
								First: value.AdminDivision{
									Code: "11",
									Name: "Foo",
								},
								Second: value.AdminDivision{
									Code: "122",
									Name: "Bar",
								},
							},
							PlaceName: "Baz",
							Position: value.Position{
								Latitude:  1.111,
								Longitude: -1.111,
							},
						},
						Distance: 0.111,
					},
					{
						PostalCode: PostalCode{
							Code:        "2222",
							CountryCode: "GB",
							AdminDivisions: value.AdminDivisions{
								First: value.AdminDivision{
									Code: "22",
									Name: "FooBar",
								},
								Second: value.AdminDivision{
									Code: "233",
									Name: "BarBaz",
								},
							},
							PlaceName: "BazFoo",
							Position: value.Position{
								Latitude:  2.222,
								Longitude: -2.222,
							},
						},
						Distance: 0.222,
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
								"/findNearbyPostalCodesJSON",
								url.Values{
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "postalcodes_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyPostalCodesRequest]{
				ctx: context.Background(),
				req: FindNearbyPostalCodesRequest{},
			},
			exp: exp[[]PostalCodeNearby]{
				res: []PostalCodeNearby{},
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
