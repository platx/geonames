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

func Test_Client_GeoCodeAddress(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req GeoCodeAddressRequest) (Address, error) {
		return client.GeoCodeAddress
	}

	testCases := []testSuite[GeoCodeAddressRequest, Address]{
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
								"/geoCodeAddressJSON",
								url.Values{
									"q":          []string{"test"},
									"country":    []string{"XX"},
									"postalcode": []string{"XXXXX"},
									"username":   []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "address_single.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[GeoCodeAddressRequest]{
				ctx: context.Background(),
				req: GeoCodeAddressRequest{
					Query:      "test",
					Country:    "XX",
					PostalCode: "XXXXX",
				},
			},
			exp: exp[Address]{
				res: Address{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					CountryCode: value.CountryCodeUnitedStates,
					AdminDivision: value.AdminDivisions{
						First: value.AdminDivision{
							ID:   0,
							Code: "D11",
							Name: "Test division 11",
						},
						Second: value.AdminDivision{
							ID:   0,
							Code: "D12",
							Name: "Test division 12",
						},
						Third: value.AdminDivision{
							ID:   0,
							Code: "D13",
							Name: "Test division 13",
						},
						Fourth: value.AdminDivision{
							ID:   0,
							Code: "D14",
							Name: "Test division 14",
						},
						Fifth: value.AdminDivision{
							ID:   0,
							Code: "D15",
							Name: "Test division 15",
						},
					},
					PostalCode:  "XXXXX",
					Locality:    "Test locality 1",
					Street:      "Test street 1",
					HouseNumber: "11A",
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
