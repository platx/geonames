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

func Test_Client_PostalCodeLookup(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req PostalCodeLookupRequest) ([]PostalCode, error) {
		return client.PostalCodeLookup
	}

	testCases := []testSuite[PostalCodeLookupRequest, []PostalCode]{
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
								"/postalCodeLookupJSON",
								url.Values{
									"postalcode": []string{"1111"},
									"maxRows":    []string{"2"},
									"country":    []string{"GB", "US"},
									"username":   []string{"test-user"},
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
			args: args[PostalCodeLookupRequest]{
				ctx: context.Background(),
				req: PostalCodeLookupRequest{
					PostalCode: "1111",
					MaxRows:    2,
					Country:    []value.CountryCode{value.CountryCodeUnitedKingdom, value.CountryCodeUnitedStates},
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
								"/postalCodeLookupJSON",
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
			args: args[PostalCodeLookupRequest]{
				ctx: context.Background(),
				req: PostalCodeLookupRequest{},
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
