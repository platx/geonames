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

func Test_Client_PostalCodeSearch(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req PostalCodeSearchRequest) ([]PostalCode, error) {
		return client.PostalCodeSearch
	}

	testCases := []testSuite[PostalCodeSearchRequest, []PostalCode]{
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
								"/postalCodeSearchJSON",
								url.Values{
									"postalcode":            []string{"1111"},
									"postalcode_startsWith": []string{"11"},
									"placename":             []string{"FooBar"},
									"placename_startsWith":  []string{"Foo"},
									"maxRows":               []string{"2"},
									"country":               []string{"GB", "US"},
									"countryBias":           []string{"GB"},
									"operator":              []string{"AND"},
									"isReduced":             []string{"true"},
									"west":                  []string{"1"},
									"east":                  []string{"2"},
									"north":                 []string{"-1"},
									"south":                 []string{"-2"},
									"username":              []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "postalcodes.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[PostalCodeSearchRequest]{
				ctx: context.Background(),
				req: PostalCodeSearchRequest{
					PostalCode:           "1111",
					PostalCodeStartsWith: "11",
					PlaceName:            "FooBar",
					PlaceNameStartsWith:  "Foo",
					MaxRows:              2,
					Country:              []value.CountryCode{value.CountryCodeUnitedKingdom, value.CountryCodeUnitedStates},
					CountryBias:          value.CountryCodeUnitedKingdom,
					Operator:             value.OperatorAnd,
					BoundingBox: value.BoundingBox{
						West:  1.0,
						East:  2.0,
						North: -1.0,
						South: -2.0,
					},
					Reduced: true,
				},
			},
			exp: exp[[]PostalCode]{
				res: []PostalCode{
					{
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
					{
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
								"/postalCodeSearchJSON",
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
			args: args[PostalCodeSearchRequest]{
				ctx: context.Background(),
				req: PostalCodeSearchRequest{},
			},
			exp: exp[[]PostalCode]{
				res: []PostalCode{},
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
