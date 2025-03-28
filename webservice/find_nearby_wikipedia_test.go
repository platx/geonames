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

func Test_Client_FindNearbyWikipedia(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req FindNearbyWikipediaRequest) ([]WikipediaNearby, error) {
		return client.FindNearbyWikipedia
	}

	testCases := []testSuite[FindNearbyWikipediaRequest, []WikipediaNearby]{
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
								"/findNearbyWikipediaJSON",
								url.Values{
									"lang":     []string{"en"},
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"radius":   []string{"10"},
									"maxRows":  []string{"2"},
									"country":  []string{"GB", "US"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "wikipedia_nearby.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyWikipediaRequest]{
				ctx: context.Background(),
				req: FindNearbyWikipediaRequest{
					Language: "en",
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Radius:  10,
					MaxRows: 2,
					Country: []value.CountryCode{value.CountryCodeUnitedKingdom, value.CountryCodeUnitedStates},
				},
			},
			exp: exp[[]WikipediaNearby]{
				res: []WikipediaNearby{
					{
						Wikipedia: Wikipedia{
							ID:          1,
							CountryCode: value.CountryCodeUnitedStates,
							Position: value.Position{
								Latitude:  1.111,
								Longitude: -1.111,
							},
							Feature:      "foo",
							Elevation:    111,
							Rank:         100,
							Language:     "en",
							WikipediaURL: "example.com/foo",
							Title:        "Foo",
							Summary:      "My to considered delightful invitation announcing of no decisively boisterous. Did add dashwoods deficient man concluded additions resources.",
						},
						Distance: 0.111,
					},
					{
						Wikipedia: Wikipedia{
							ID:          2,
							CountryCode: value.CountryCodeUnitedKingdom,
							Position: value.Position{
								Latitude:  2.222,
								Longitude: -2.222,
							},
							Feature:      "bar",
							Elevation:    222,
							Rank:         200,
							Language:     "es",
							WikipediaURL: "example.com/bar",
							Title:        "Bar",
							Summary:      "Full he none no side. Uncommonly surrounded considered for him are its. It we is read good soon.",
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
								"/findNearbyWikipediaJSON",
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
			args: args[FindNearbyWikipediaRequest]{
				ctx: context.Background(),
				req: FindNearbyWikipediaRequest{},
			},
			exp: exp[[]WikipediaNearby]{
				res: []WikipediaNearby{},
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
