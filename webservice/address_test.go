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

func Test_Client_Address(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req AddressRequest) ([]AddressNearby, error) {
		return client.Address
	}

	testCases := []testSuite[AddressRequest, []AddressNearby]{
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
								"/addressJSON",
								url.Values{
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"radius":   []string{"11"},
									"maxRows":  []string{"2"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "address_nearby.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[AddressRequest]{
				ctx: context.Background(),
				req: AddressRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Radius:  11,
					MaxRows: 2,
				},
			},
			exp: exp[[]AddressNearby]{
				res: []AddressNearby{
					{
						Address: Address{
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
							SourceID:    "111",
							PostalCode:  "XXXXX",
							Locality:    "Test locality 1",
							Street:      "Test street 1",
							HouseNumber: "11A",
						},
						Distance: 0.111,
					},
					{
						Address: Address{
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
							SourceID:    "222",
							PostalCode:  "YYYYY",
							Locality:    "Test locality 2",
							Street:      "Test street 2",
							HouseNumber: "22B",
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
								"/addressJSON",
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
			args: args[AddressRequest]{
				ctx: context.Background(),
				req: AddressRequest{},
			},
			exp: exp[[]AddressNearby]{
				res: []AddressNearby{},
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
			args: args[AddressRequest]{
				ctx: context.Background(),
				req: AddressRequest{},
			},
			exp: exp[[]AddressNearby]{
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
			args: args[AddressRequest]{
				ctx: context.Background(),
				req: AddressRequest{},
			},
			exp: exp[[]AddressNearby]{
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
			args: args[AddressRequest]{
				ctx: context.Background(),
				req: AddressRequest{},
			},
			exp: exp[[]AddressNearby]{
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
			args: args[AddressRequest]{
				ctx: context.Background(),
				req: AddressRequest{},
			},
			exp: exp[[]AddressNearby]{
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
			args: args[AddressRequest]{
				ctx: nil,
				req: AddressRequest{},
			},
			exp: exp[[]AddressNearby]{
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
