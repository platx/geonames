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
									"type":         []string{"json"},
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
							FeatureClass:     "A",
							FeatureClassName: "Test class",
							FeatureCode:      "AAAA",
							FeatureCodeName:  "Test code",
							Name:             "New York City",
							ToponymName:      "New York City",
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
							FeatureClass:     "A",
							FeatureClassName: "Test class",
							FeatureCode:      "AAAA",
							FeatureCodeName:  "Test code",
							Name:             "London",
							ToponymName:      "London",
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
									"type":     []string{"json"},
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
		{
			name: "invalid country id",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"countryId": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => parse CountryID => strconv.ParseUint: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid distance",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"distance": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => parse Distance => strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid longitude",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"lng": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => parse Position => longitude => strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid latitude",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"lat": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("decode response => parse Position => latitude => strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
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
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
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
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
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
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
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
			args: args[FindNearbyRequest]{
				ctx: context.Background(),
				req: FindNearbyRequest{},
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
			args: args[FindNearbyRequest]{
				ctx: nil,
				req: FindNearbyRequest{},
			},
			exp: exp[[]GeoNameNearby]{
				res: nil,
				err: errors.New("create http request => net/http: nil Context"),
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}
