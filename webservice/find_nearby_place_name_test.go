package webservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_FindNearbyPlaceName(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req FindNearbyPlaceNameRequest) ([]GeoNameNearby, error) {
		return client.FindNearbyPlaceName
	}

	testCases := []testSuite[FindNearbyPlaceNameRequest, []GeoNameNearby]{
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
								"/findNearbyPlaceNameJSON",
								url.Values{
									"lat":          []string{"1.111"},
									"lng":          []string{"-1.111"},
									"lang":         []string{"en"},
									"radius":       []string{"10"},
									"maxRows":      []string{"2"},
									"localCountry": []string{"true"},
									"cities":       []string{"cities5000"},
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
			args: args[FindNearbyPlaceNameRequest]{
				ctx: context.Background(),
				req: FindNearbyPlaceNameRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Language:     "en",
					Radius:       10,
					MaxRows:      2,
					LocalCountry: true,
					Cities:       value.Cities5000,
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
								"/findNearbyPlaceNameJSON",
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
			args: args[FindNearbyPlaceNameRequest]{
				ctx: context.Background(),
				req: FindNearbyPlaceNameRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: []GeoNameNearby{},
				err: nil,
			},
		},
		{
			name: "invalid success response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geo`)),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyPlaceNameRequest]{
				ctx: context.Background(),
				req: FindNearbyPlaceNameRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "error response",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       testutil.MustOpen(testdata.FS, "authorization_error.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyPlaceNameRequest]{
				ctx: context.Background(),
				req: FindNearbyPlaceNameRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => got error response => code: 10, message: \"user does not exist.\""),
			},
		},
		{
			name: "invalid error response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(strings.NewReader(`{"stat`)),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyPlaceNameRequest]{
				ctx: context.Background(),
				req: FindNearbyPlaceNameRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "send request failed",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				userName: "test-user",
			},
			args: args[FindNearbyPlaceNameRequest]{
				ctx: context.Background(),
				req: FindNearbyPlaceNameRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[FindNearbyPlaceNameRequest]{
				ctx: nil,
				req: FindNearbyPlaceNameRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("create http request => net/http: nil Context"),
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
