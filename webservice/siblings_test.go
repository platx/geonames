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

func Test_Client_Siblings(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req SiblingsRequest) ([]GeoName, error) {
		return client.Siblings
	}

	testCases := []testSuite[SiblingsRequest, []GeoName]{
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
								"/siblingsJSON",
								url.Values{
									"geonameId": []string{"1"},
									"type":      []string{"json"},
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{
					ID: 1,
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
				},
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
				res: nil,
				err: errors.New("decode response => parse CountryID => strconv.ParseUint: parsing \"invalid\": invalid syntax"),
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
			args: args[SiblingsRequest]{
				ctx: context.Background(),
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[SiblingsRequest]{
				ctx: nil,
				req: SiblingsRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
