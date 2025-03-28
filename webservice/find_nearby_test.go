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

func Test_Client_FindNearby(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req FindNearbyRequest) ([]GeoNameNearby, error) {
		return client.FindNearby
	}

	testCases := []testSuite[FindNearbyRequest, []GeoNameNearby]{
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
								"/findNearbyJSON",
								url.Values{
									"lat":          []string{"1.111"},
									"lng":          []string{"-1.111"},
									"featureClass": []string{"A"},
									"featureCode":  []string{"AAAA", "AAAB"},
									"radius":       []string{"10"},
									"maxRows":      []string{"2"},
									"localCountry": []string{"true"},
									"username":     []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geonames_nearby.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					FeatureClass: []string{"A"},
					FeatureCode:  []string{"AAAA", "AAAB"},
					Radius:       10,
					MaxRows:      2,
					LocalCountry: true,
				},
			},
			exp: exp[[]GeoNameNearby]{
				res: []GeoNameNearby{
					{
						GeoName: GeoName{
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
								Second: value.AdminDivision{},
								Third:  value.AdminDivision{},
								Fourth: value.AdminDivision{},
								Fifth:  value.AdminDivision{},
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
						Distance: 0.111,
					},
					{
						GeoName: GeoName{
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
								Second: value.AdminDivision{},
								Third:  value.AdminDivision{},
								Fourth: value.AdminDivision{},
								Fifth:  value.AdminDivision{},
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
								"/findNearbyJSON",
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
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: []GeoNameNearby{},
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
