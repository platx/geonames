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
			args: args[PostalCodeLookupRequest]{
				ctx: context.Background(),
				req: PostalCodeLookupRequest{},
			},
			exp: exp[[]PostalCode]{
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
			args: args[PostalCodeLookupRequest]{
				ctx: context.Background(),
				req: PostalCodeLookupRequest{},
			},
			exp: exp[[]PostalCode]{
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
			args: args[PostalCodeLookupRequest]{
				ctx: context.Background(),
				req: PostalCodeLookupRequest{},
			},
			exp: exp[[]PostalCode]{
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
			args: args[PostalCodeLookupRequest]{
				ctx: context.Background(),
				req: PostalCodeLookupRequest{},
			},
			exp: exp[[]PostalCode]{
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
			args: args[PostalCodeLookupRequest]{
				ctx: nil,
				req: PostalCodeLookupRequest{},
			},
			exp: exp[[]PostalCode]{
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
