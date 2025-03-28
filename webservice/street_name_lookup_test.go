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

func Test_Client_StreetNameLookup(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req StreetNameLookupRequest) ([]Address, error) {
		return client.StreetNameLookup
	}

	testCases := []testSuite[StreetNameLookupRequest, []Address]{
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
								"/streetNameLookupJSON",
								url.Values{
									"q":                  []string{"test"},
									"country":            []string{"XX"},
									"postalcode":         []string{"XXXXX"},
									"adminCode1":         []string{"D1"},
									"adminCode2":         []string{"D2"},
									"adminCode3":         []string{"D3"},
									"isUniqueStreetName": []string{"true"},
									"username":           []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "address.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[StreetNameLookupRequest]{
				ctx: context.Background(),
				req: StreetNameLookupRequest{
					Query:      "test",
					Country:    "XX",
					PostalCode: "XXXXX",
					AdminCode: value.AdminCode{
						First:  "D1",
						Second: "D2",
						Third:  "D3",
					},
					UniqueStreetName: true,
				},
			},
			exp: exp[[]Address]{
				res: []Address{
					{
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
					{
						Position: value.Position{
							Latitude:  2.222,
							Longitude: -2.222,
						},
						CountryCode: value.CountryCodeUnitedKingdom,
						AdminDivision: value.AdminDivisions{
							First: value.AdminDivision{
								ID:   0,
								Code: "D21",
								Name: "Test division 21",
							},
							Second: value.AdminDivision{
								ID:   0,
								Code: "D22",
								Name: "Test division 22",
							},
							Third: value.AdminDivision{
								ID:   0,
								Code: "D23",
								Name: "Test division 23",
							},
							Fourth: value.AdminDivision{
								ID:   0,
								Code: "D24",
								Name: "Test division 24",
							},
							Fifth: value.AdminDivision{
								ID:   0,
								Code: "D25",
								Name: "Test division 25",
							},
						},
						PostalCode:  "YYYYY",
						Locality:    "Test locality 2",
						Street:      "Test street 2",
						HouseNumber: "22B",
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
								"/streetNameLookupJSON",
								url.Values{
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "address_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[StreetNameLookupRequest]{
				ctx: context.Background(),
				req: StreetNameLookupRequest{},
			},
			exp: exp[[]Address]{
				res: []Address{},
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
